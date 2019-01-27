package target

import (
	"fmt"
	"io/ioutil"

	"github.com/ssoor/kuberes/pkg/loader"
	"github.com/ssoor/kuberes/pkg/reference"
	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
)

const (
	loaderPathTarget        = "kuberes.yaml"
	configPathReferenceRule = ".kuberes/reference.yaml"
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

	decoder, err = t.loader.LoadYamlDecoder(configPathReferenceRule)
	if nil != err {
		body, err := ioutil.ReadFile(configPathReferenceRule)
		if nil != err {
			return err
		}

		decoder = yaml.NewFormatErrorDecodeFromBytes(body, configPathReferenceRule)
	}

	if err := t.referenceMap.Load(decoder); nil != err {
		return err
	}

	return nil
}

// Make is
func (t *Target) Make() (ResourceMap, error) {
	if err := t.loadImports(); nil != err {
		return nil, err
	}

	if err := t.loadResources(); nil != err {
		return nil, err
	}

	if err := t.Patchs.Make(t.loader, t.resourceMap); nil != err {
		return nil, err
	}

	err := t.resourceMap.Range(func(id resource.UniqueID, res *resource.Resource) error {
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
		return nil, err
	}

	return t.generateResourceMap()
}

func (t *Target) loadImports() error {
	for _, depend := range t.Imports {
		resourceMap, err := depend.Make(t.loader)
		if nil != err {
			return err
		}

		if err := t.resourceMap.Merge(false, resourceMap); nil != err {
			return err
		}
	}

	return nil
}

func (t *Target) loadResources() error {
	for _, path := range t.Resources {
		decoder, err := t.loader.LoadYamlDecoder(path)
		if nil != err {
			return err
		}

		if err := t.resourceMap.MergeFormDecoder(false, decoder); nil != err {
			return err
		}
	}

	return nil
}

func (t *Target) generateResourceMap() (ResourceMap, error) {
	resultMap := make(ResourceMap)

	err := t.resourceMap.Range(func(id resource.UniqueID, res *resource.Resource) error {
		if rule := t.referenceMap.FindByGVK(res.GVK()); nil != rule {
			if err := rule.Refresh(res, t.referenceCallback); nil != err {
				return err
			}
		}

		resultMap[res.ID()] = res

		return nil
	})

	if nil != err {
		return nil, err
	}

	return resultMap, nil
}

func (t *Target) referenceCallback(fs reference.RefreshSpec, in interface{}) (out interface{}, err error) {
	out = in

	switch fs.Name {
	case "matedata.name":
		id := resource.NewUniqueID(in.(string), fs.GVK)

		res := t.resourceMap[id]
		if nil == res {
			panic("TODO")
		}

		out = res.GetName()
	case "matedata.labels":
		out = t.Matedata.Labels
	case "matedata.annotations":
		out = t.Matedata.Annotations
	}

	fmt.Printf("[%s] %s:%s(%v) => %v\n", fs.Resource.ID(), fs.Name, fs.Path, in, out)

	return out, nil
}
