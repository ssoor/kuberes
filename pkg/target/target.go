package target

import (
	"fmt"

	"github.com/ssoor/kuberes/pkg/loader"
	"github.com/ssoor/kuberes/pkg/reference"
	"github.com/ssoor/kuberes/pkg/resource"
)

const (
	loaderPathTarget      = "kuberes.yaml"
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

	loader       loader.Loader
	resourceMap  ResourceMap
	referenceMap ReferenceMap
}

// ResourceMap is
func (t *Target) ResourceMap() ResourceMap {
	return t.resourceMap
}

// NewTarget is
func NewTarget(loader loader.Loader) (*Target, error) {
	newTarget := &Target{
		loader:       loader,
		resourceMap:  make(ResourceMap),
		referenceMap: make(ReferenceMap),
	}

	return newTarget, nil
}

// Load is
func (t *Target) Load() (err error) {
	decoder, err := t.loader.LoadYamlDecoder(loaderPathTarget)
	if nil != err {
		return err
	}

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
	for _, depend := range t.Imports {
		resourceMap, err := depend.Make(t.loader)
		if nil != err {
			return err
		}

		if err := t.resourceMap.Merge(false, resourceMap); nil != err {
			return err
		}
	}

	for _, path := range t.Resources {
		decoder, err := t.loader.LoadYamlDecoder(path)
		if nil != err {
			return err
		}

		if err := t.resourceMap.MergeFormDecoder(false, decoder); nil != err {
			return err
		}
	}

	t.Patchs.Make(t.loader, t.resourceMap)

	for key := range t.referenceMap {
		fmt.Println(key)
	}

	err = t.resourceMap.Range(func(id resource.UniqueID, res *resource.Resource) error {
		switch name := res.GetName(); name {
		case "":
		case "-":
			res.SetName(t.Name)
		default:
			res.SetName(fmt.Sprintf("%s-%s", t.Name, res.GetName()))
		}

		return t.Matedata.Make(res)
	})
	if nil != err {
		return err
	}

	err = t.resourceMap.Range(func(id resource.UniqueID, res *resource.Resource) error {
		fmt.Printf("id => %v ,res.ID() => %v\n", id, res.ID())

		refRule := t.referenceMap.FindByGVK(res.GVK())
		if nil == refRule {
			return nil // continue
		}

		err = refRule.RefreshMatedataName(res, func(fs reference.FieldSpec, fp reference.FieldPath, in interface{}) (interface{}, error) {
			id := resource.NewUniqueID(in.(string), res.GetNamespace(), fs.GVK)

			res := t.resourceMap[id]

			fmt.Println(id, fp)
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

		return err
	})

	return err
}
