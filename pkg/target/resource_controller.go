package target

import (
	"fmt"
	"io"

	"sigs.k8s.io/kustomize/pkg/transformers/config/defaultconfig"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
)

// ResourceController is
type ResourceController struct {
	resources ResourceMap
}

// NewResourceController is
func NewResourceController() (*ResourceController, error) {
	control := &ResourceController{
		resources: make(ResourceMap),
	}

	return control, nil
}

// Map is
func (rc *ResourceController) Map() ResourceMap {
	return rc.resources
}

// Range is
func (rc *ResourceController) Range(fn func(resource.UniqueID, *resource.Resource) error) error {
	for id, res := range rc.resources {
		if err := fn(id, res); nil != err {
			return err
		}
	}

	return nil
}

// LoadResourcesFormBytes is
func (rc *ResourceController) LoadResourcesFormBytes(decoder yaml.Decoder) ([]*resource.Resource, error) {
	var err error
	var result []*resource.Resource

	defaultconfig.GetDefaultFieldSpecs()

	for err == nil {
		var out unstructured.Unstructured

		if err = decoder.Decode(&out); err != nil {
			continue
		}

		res := &resource.Resource{Unstructured: out}

		if err := res.Validate(); err != nil {
			return nil, err
		}

		result = append(result, res)
	}

	if err != io.EOF {
		return nil, err
	}

	return result, nil
}

// MergeResourceMap is
func (rc *ResourceController) MergeResourceMap(override bool, resources ...ResourceMap) error {
	for _, m := range resources {
		if m == nil {
			continue
		}

		for _, res := range m {
			if err := rc.MergeResources(override, res); nil != err {
				return err
			}
		}
	}

	return nil
}

// MergeResources is
func (rc *ResourceController) MergeResources(override bool, resources ...*resource.Resource) error {
	for _, res := range resources {
		id := res.ID()

		if !override {
			if _, found := rc.resources[id]; found {
				return fmt.Errorf("id '%q' already used", id)
			}
		}

		rc.resources[id] = res
	}

	return nil
}

// Find is
func (rc *ResourceController) Find(id resource.UniqueID) *resource.Resource {
	for id, res := range rc.resources {
		if id.Equals(id) {
			return res
		}
	}

	return nil
}
