package target

import (
	"fmt"

	"github.com/ssoor/kuberes/pkg/gvk"
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

// Patchs is
type Patchs struct {
	RFC6902   []RFC6902Patch `json:"rfc6902,omitempty" yaml:"rfc6902,omitempty"`
	Strategic []string       `json:"strategic,omitempty" yaml:"strategic,omitempty"`
}

// Make is
func (p Patchs) Make(loader loader.Loader, resourceMap ResourceMap) error {
	for _, patch := range p.RFC6902 {
		body, err := loader.LoadBytes(patch.Path)
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

		merger := merge.NewRFC6902Merger(body)

		p.resourceMerge(patch.UniqueID, merger, resourceMap)
	}

	for _, path := range p.Strategic {
		decoder, err := loader.LoadYamlDecoder(path)
		if nil != err {
			return err
		}

		patchs, err := resource.NewResourceFromDecoder(decoder)
		if nil != err {
			return err
		}

		for _, patch := range patchs {
			_, err := scheme.Scheme.New(p.toSchemaGvk(patch.GVK()))

			var merger merge.Merger
			switch {
			case runtime.IsNotRegisteredError(err):
				merger = merge.NewJSONMerger(patch.Map())
			case nil == err:
				merger = merge.NewStrategicMerge(patch.Map())
			default:
				return err
			}

			p.resourceMerge(patch.ID(), merger, resourceMap)
		}
	}

	return nil
}

func (p Patchs) resourceMerge(uniqueID resource.UniqueID, merger merge.Merger, resourceMap ResourceMap) error {
	res, exists := resourceMap[uniqueID]
	if !exists {
		return fmt.Errorf("failed to find an object with %s to apply the patch", uniqueID)
	}

	versionObj, err := scheme.Scheme.New(p.toSchemaGvk(uniqueID.GVK))
	if nil != err {
		versionObj = nil
	}

	resMap, err := merger.MergeMap(res.Map(), versionObj)
	if nil != err {
		return err
	}

	res.SetMap(map[string]interface{}(resMap))

	return nil
}

// toSchemaGvk converts to a schema.GroupVersionKind.
func (p Patchs) toSchemaGvk(x gvk.GVK) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   x.Group,
		Version: x.Version,
		Kind:    x.Kind,
	}
}
