package target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"sigs.k8s.io/kustomize/pkg/constants"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/types"

	"github.com/pkg/errors"

	"github.com/ssoor/kuberes/pkg/log"
	"github.com/ssoor/kuberes/pkg/reference"
	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
)

const (
	fileNameReferenceRule = "conf/reference_rule.yaml"
)

// Target is
type Target struct {
	ldr           ifc.Loader
	rules         ReferenceRuleMap
	resources     *ResourceController
	kustomization *types.Kustomization
}

// Resources is
func (t *Target) Resources() *ResourceController {
	return t.resources
}

// NewTarget is
func NewTarget(ldr ifc.Loader) (*Target, error) {
	content, err := loadKustFile(ldr)
	if err != nil {
		return nil, err
	}

	var k types.Kustomization
	if err := unmarshal(content, &k); err != nil {
		return nil, err
	}

	k.DealWithDeprecatedFields()
	msgs, errs := k.EnforceFields()
	if len(errs) > 0 {
		return nil, fmt.Errorf(strings.Join(errs, "\n"))
	}
	if len(msgs) > 0 {
		log.Printf(strings.Join(msgs, "\n"))
	}

	newTarget := &Target{
		ldr:           ldr,
		kustomization: &k,
	}

	newTarget.rules, err = LoadReferenceRuleMapFormFile(fileNameReferenceRule)
	if nil != err {
		return nil, err
	}

	return newTarget, nil
}

func loadKustFile(ldr ifc.Loader) ([]byte, error) {
	for _, kf := range []string{
		constants.KustomizationFileName,
		constants.SecondaryKustomizationFileName} {
		content, err := ldr.Load(kf)
		if err == nil {
			return content, nil
		}
		if !strings.Contains(err.Error(), "no such file or directory") {
			return nil, err
		}
	}
	return nil, fmt.Errorf("no kustomization.yaml file under %s", ldr.Root())
}

func unmarshal(y []byte, o interface{}) error {
	j, err := yaml.ToJSONFormBytes(y)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(j))
	dec.DisallowUnknownFields()
	return dec.Decode(o)
}

// Load is
func (t *Target) Load() (err error) {
	if t.resources, err = NewResourceController(); nil != err {
		return err
	}

	for _, path := range t.kustomization.Resources {
		content, err := t.ldr.Load(path)
		if err != nil {
			return errors.Wrap(err, "Load from path "+path+" failed")
		}

		res, err := t.resources.LoadResourcesFormBytes(yaml.NewFormatErrorDecodeFormBytes(content, path))
		if err != nil {
			return err
		}

		if err := t.resources.MergeResources(false, res...); nil != err {
			return err
		}
	}

	return nil
}

// RefreshReferences is
func (t *Target) RefreshReferences() (err error) {
	for key := range t.rules {
		fmt.Println(key)
	}

	t.resources.Range(func(id resource.UniqueID, res *resource.Resource) error {
		fmt.Printf("id => %v ,res.ID() => %v\n", id, res.ID())

		references, exists := t.rules[res.GVKID()]
		if !exists {
			return nil // continue
		}

		err = t.refreshFields(res, references.MatedataName, func(fs reference.FieldSpec, fp reference.FieldPath, in interface{}) (interface{}, error) {
			id := resource.NewUniqueID(res.GetName(), res.GetNamespace(), fs.GVK)
			res := t.resources.Find(id)

			fmt.Println(id, fp, res.GetName())
			return res.GetName(), nil
		})
		if nil != err {
			return err
		}

		err = t.refreshFields(res, references.MatedataLabels, func(fs reference.FieldSpec, fp reference.FieldPath, in interface{}) (interface{}, error) {
			return res.GetLabels(), nil
		})
		if nil != err {
			return err
		}

		err = t.refreshFields(res, references.MatedataAnnotations, func(fs reference.FieldSpec, fp reference.FieldPath, in interface{}) (interface{}, error) {
			return res.GetAnnotations(), nil
		})
		if nil != err {
			return err
		}

		return nil
	})

	return nil
}

func (t Target) refreshFields(res *resource.Resource, fields []reference.FieldSpec, fn func(reference.FieldSpec, reference.FieldPath, interface{}) (interface{}, error)) error {
	for _, field := range fields {
		if err := field.Refresh(res.Map(), fn); nil != err {
			return err
		}
	}

	return nil
}
