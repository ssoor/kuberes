package target

import (
	"fmt"

	"github.com/ssoor/kuberes/pkg/loader"
	"github.com/ssoor/kuberes/pkg/resource"
)

const (
	loaderPathTarget        = "kuberes.yaml"
	configPathReferenceRule = ".kuberes/reference.yaml"
)

type Patchs struct {
	RFC6902   []RFC6902Patch `json:"rfc6902,omitempty" yaml:"rfc6902,omitempty"`
	Strategic []string       `json:"strategic,omitempty" yaml:"strategic,omitempty"`
}

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
}

// Maker is
type Maker interface {
	Make() (map[resource.UniqueID]*resource.Resource, error)
}

// targetMake is
type targetMake struct {
	name   string
	conf   *Target
	loader loader.Loader
	resc   ResourceController
	refc   *ReferenceControl
	patchc PatchController
}

// NewMaker is
func NewMaker(loader loader.Loader, name string) (Maker, error) {
	conf := &Target{}

	decoder, err := loader.LoadYamlDecoder(loaderPathTarget)
	if nil != err {
		return nil, err
	}

	if err := decoder.Decode(conf); nil != err {
		return nil, err
	}

	resc, err := NewResourceControl(loader, conf.Resources)
	if nil != err {
		return nil, nil
	}

	refc, err := NewReferenceControl(loader, loaderPathTarget)
	if nil != err {
		return nil, nil
	}

	patchc, err := NewPatchController(loader, conf.Patchs.RFC6902, conf.Patchs.Strategic)
	if nil != err {
		return nil, nil
	}

	t := &targetMake{
		name:   name,
		conf:   conf,
		resc:   resc,
		refc:   refc,
		patchc: patchc,
	}

	refc.AddRefreshHandle("matedata.name", t.refreshName)
	refc.AddRefreshHandle("matedata.labels", t.refreshLabels)
	refc.AddRefreshHandle("matedata.annotations", t.refreshAnnotations)

	return t, nil
}

// Make is
func (t *targetMake) Make() (resourceMap map[resource.UniqueID]*resource.Resource, err error) {
	if err := t.loadImports(); nil != err {
		return nil, err
	}

	if err := t.resc.Range(func(id resource.UniqueID, res *resource.Resource) error {
		if err := t.patchc.Patch(res); nil != err {
			return err
		}

		if err := t.modify(res); nil != err {
			return err
		}

		return nil
	}); nil != err {
		return nil, err
	}

	resourceMap = make(map[resource.UniqueID]*resource.Resource)
	if err := t.resc.Range(func(id resource.UniqueID, res *resource.Resource) error {
		if err := t.refc.Refresh(res); nil != err {
			return err
		}

		resourceMap[res.ID()] = res

		return nil
	}); nil != err {
		return nil, err
	}

	return resourceMap, nil
}

func (t targetMake) modify(res *resource.Resource) error {
	targetName := t.name

	switch name := res.GetName(); name {
	case "":
	case "-":
		res.SetName(targetName)
	default:
		res.SetName(fmt.Sprintf("%s-%s", targetName, res.GetName()))
	}

	m := t.conf.Matedata

	res.SetNamespace(m.Namespace)
	res.SetLabels(t.mergeMap(res.GetLabels(), m.Labels))
	res.SetAnnotations(t.mergeMap(res.GetAnnotations(), m.Annotations))

	return nil

}

func (t *targetMake) loadImports() error {
	for _, depend := range t.conf.Imports {
		resourceMap, err := depend.Make(t.loader)
		if nil != err {
			return err
		}

		for id, res := range resourceMap {
			if err := t.resc.Add(id, res, false); nil != err {
				return err
			}
		}

	}

	return nil
}

// MergeMap is
func (t targetMake) mergeMap(src map[string]string, merge map[string]string) map[string]string {
	if nil == merge {
		return src
	}

	if nil == src {
		return merge
	}

	for key, val := range merge {
		src[key] = val
	}

	return src
}

func (t *targetMake) refreshName(fs FieldSpec, path resource.Path, in interface{}) (interface{}, error) {
	switch in.(type) {
	case string:
	default:
		panic("TODO")
	}

	id := resource.NewUniqueID(in.(string), fs.GVK)

	res := t.resc.Get(id)
	if nil == res {
		panic("TODO")
	}

	return res.GetName(), nil
}

func (t *targetMake) refreshLabels(fs FieldSpec, path resource.Path, in interface{}) (interface{}, error) {
	var val map[string]string

	switch inVal := in.(type) {
	case map[string]string:
		val = inVal
	}

	return t.mergeMap(val, t.conf.Matedata.Labels), nil
}

func (t *targetMake) refreshAnnotations(fs FieldSpec, path resource.Path, in interface{}) (interface{}, error) {
	var val map[string]string

	switch inVal := in.(type) {
	case map[string]string:
		val = inVal
	}

	return t.mergeMap(val, t.conf.Matedata.Annotations), nil
}
