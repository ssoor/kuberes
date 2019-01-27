package target

import (
	"fmt"
	"io/ioutil"

	"github.com/ssoor/kuberes/pkg/reference"
	"github.com/ssoor/kuberes/pkg/resource"
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

	resourceMap  ResourceMap
	referenceMap ReferenceMap
}

// ResourceMap is
func (t *Target) ResourceMap() ResourceMap {
	return t.resourceMap
}

// NewTarget is
func NewTarget() (*Target, error) {
	newTarget := &Target{
		resourceMap:  make(ResourceMap),
		referenceMap: make(ReferenceMap),
	}

	return newTarget, nil
}

// Load is
func (t *Target) Load(path string) (err error) {
	body, err := ioutil.ReadFile(path)
	if nil != err {
		return err
	}

	decoder := yaml.NewFormatErrorDecodeFromBytes(body, path)
	if err := decoder.Decode(t); nil != err {
		return err
	}

	if t.referenceMap.Load(fileNameReferenceRule); nil != err {
		return err
	}

	return nil
}

// Make is
func (t *Target) Make() (err error) {
	for key := range t.referenceMap {
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

		refRule := t.referenceMap.FindByGVK(res.GVKID())
		if nil == refRule {
			return nil // continue
		}

		err = refRule.RefreshMatedataName(res, func(fs reference.FieldSpec, fp reference.FieldPath, in interface{}) (interface{}, error) {
			id := resource.NewUniqueID(res.GetName(), res.GetNamespace(), fs.GVK)
			res := t.resourceMap[id]

			fmt.Println(id, fp, res.GetName())
			return res.GetName(), nil
		})
		if nil != err {
			return err
		}

		err = refRule.RefreshMatedataLabels(res, func(fs reference.FieldSpec, fp reference.FieldPath, in interface{}) (interface{}, error) {
			return res.GetLabels(), nil
		})
		if nil != err {
			return err
		}

		err = refRule.RefreshMatedataAnnotations(res, func(fs reference.FieldSpec, fp reference.FieldPath, in interface{}) (interface{}, error) {
			return res.GetAnnotations(), nil
		})
		if nil != err {
			return err
		}

		return nil
	})

	return nil
}
