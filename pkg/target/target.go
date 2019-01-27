package target

import (
	"fmt"
	"io/ioutil"

	"github.com/ssoor/kuberes/pkg/reference"
	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/resourcemap"
	"github.com/ssoor/kuberes/pkg/yaml"
)

const (
	fileNameReferenceRule = "conf/reference_rule.yaml"
)

// Target is
type Target struct {
	Name    string   `json:"name,omitempty" yaml:"name,omitempty"`
	Imports []Import `json:"imports,omitempty" yaml:"imports,omitempty"`

	// Patchs to add to all objects.
	Patchs Patchs `json:"patchs,omitempty" yaml:"patchs,omitempty"`

	// Matedata to add to all objects.
	Matedata Matedata `json:"matedata,omitempty" yaml:"matedata,omitempty"`

	// Resources specifies relative paths to files holding YAML representations
	// of kubernetes API objects. URLs and globs not supported.
	Resources []string `json:"resources,omitempty" yaml:"resources,omitempty"`

	rules       ReferenceRuleMap
	resourceMap resourcemap.ResourceMap
}

// ResourceMap is
func (t *Target) ResourceMap() resourcemap.ResourceMap {
	return t.resourceMap
}

// NewTarget is
func NewTarget() (*Target, error) {
	newTarget := &Target{
		rules:       make(ReferenceRuleMap),
		resourceMap: resourcemap.New(),
	}

	return newTarget, nil
}

// Load is
func (t *Target) Load(path string) (err error) {
	body, err := ioutil.ReadFile(path)
	if nil != err {
		return err
	}

	decoder := yaml.NewFormatErrorDecodeFormBytes(body, path)
	if err := decoder.Decode(t); nil != err {
		return err
	}

	if t.rules.Load(fileNameReferenceRule); nil != err {
		return err
	}

	return nil
}

// Make is
func (t *Target) Make() (err error) {
	for key := range t.rules {
		fmt.Println(key)
	}

	for _, depend := range t.Imports {
		resourceMap, err := depend.Make()
		if nil != err {
			return err
		}

		if err := t.resourceMap.Merge(false, resourceMap); nil != err {
			return err
		}
	}

	for _, path := range t.Resources {
		if err := t.resourceMap.MergeFormPath(false, path); nil != err {
			return err
		}
	}

	t.resourceMap.Range(func(id resource.UniqueID, res *resource.Resource) error {
		fmt.Printf("id => %v ,res.ID() => %v\n", id, res.ID())

		res.SetName(fmt.Sprintf("%s-%s", t.Name, res.GetName()))

		t.Patchs.MakeResource(res)
		t.Matedata.MakeResource(res)

		references, exists := t.rules[res.GVKID()]
		if !exists {
			return nil // continue
		}

		err = t.refreshFields(res, references.MatedataName, func(fs reference.FieldSpec, fp reference.FieldPath, in interface{}) (interface{}, error) {
			id := resource.NewUniqueID(res.GetName(), res.GetNamespace(), fs.GVK)
			res := t.resourceMap[id]

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
