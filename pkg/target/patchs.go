package target

import (
	"github.com/ssoor/kuberes/pkg/loader"
	"github.com/ssoor/kuberes/pkg/merge"
	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

// RFC6902Patch is
type RFC6902Patch struct {
	resource.UniqueID `json:",inline,omitempty" yaml:",inline,omitempty"`

	// relative file path for a json patch file inside a kustomization
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

// PatchController is
type PatchController interface {
	Patch(*resource.Resource) error
}

type patchControl struct {
	l    loader.Loader
	resc ResourceController

	mergerMap map[resource.UniqueID][]merge.Merger
}

// NewPatchController is
func NewPatchController(loader loader.Loader, rfc6902Patchs []RFC6902Patch, strategicPatchs []string) (PatchController, error) {
	pc := &patchControl{}

	for _, patch := range rfc6902Patchs {
		if err := pc.addRFC6902Patch(patch); nil != err {
			return nil, err
		}
	}

	for _, patch := range strategicPatchs {
		if err := pc.addStrategicPatch(patch); nil != err {
			return nil, err
		}
	}

	return pc, nil
}

// Make is
func (pc *patchControl) Patch(res *resource.Resource) (err error) {
	resMap := res.Map()
	mergers := pc.mergerMap[res.ID()]
	for _, merger := range mergers {
		versionObj, err := scheme.Scheme.New(pc.toSchemaGvk(res.GVK()))
		if nil != err {
			versionObj = nil
		}

		if resMap, err = merger.MergeMap(resMap, versionObj); nil != err {
			return err
		}
	}

	res.SetMap(map[string]interface{}(resMap))

	return nil
}

func (pc *patchControl) addRFC6902Patch(patch RFC6902Patch) error {
	body, err := pc.l.LoadBytes(patch.Path)
	if nil != err {
		return err
	}

	if body[0] != '[' {
		// if it isn't JSON, try to parse it as YAML
		body, err = yaml.ToJSON(body)
		if err != nil {
			return err
		}
	}

	pc.mergerMap[patch.UniqueID] = append(pc.mergerMap[patch.UniqueID], merge.NewRFC6902Merger(body))

	return nil
}

func (pc *patchControl) addStrategicPatch(path string) error {
	decoder, err := pc.l.LoadYamlDecoder(path)
	if nil != err {
		return err
	}

	patchs, err := resource.NewResourceFromDecoder(decoder)
	if nil != err {
		return err
	}

	for _, patch := range patchs {
		_, err := scheme.Scheme.New(pc.toSchemaGvk(patch.GVK()))

		var merger merge.Merger
		switch {
		case runtime.IsNotRegisteredError(err):
			merger = merge.NewJSONMerger(patch.Map())
		case nil == err:
			merger = merge.NewStrategicMerge(patch.Map())
		default:
			return err
		}

		pc.mergerMap[patch.ID()] = append(pc.mergerMap[patch.ID()], merger)
	}

	return nil
}

// toSchemaGvk converts to a schema.GroupVersionKind.
func (pc patchControl) toSchemaGvk(x resource.GVK) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   x.Group,
		Version: x.Version,
		Kind:    x.Kind,
	}
}
