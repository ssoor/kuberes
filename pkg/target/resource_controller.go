package target

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"sigs.k8s.io/kustomize/pkg/transformers/config/defaultconfig"

	yaml2 "github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/ssoor/kuberes/pkg/log"
	"github.com/ssoor/kuberes/pkg/resource"
)

func isEmptyYamlError(err error) bool {
	return strings.Contains(err.Error(), "is missing in 'null'")
}

// RuleMap is a map from resource ID to Resource.
type RuleMap map[resource.ID]RuleSpec

// ResourceMap is a map from resource ID to Resource.
type ResourceMap map[resource.ID]*resource.Resource

// ResourceController is
type ResourceController struct {
	rules     RuleMap
	resources ResourceMap
}

// NewResourceController is
func NewResourceController() (*ResourceController, error) {
	control := &ResourceController{
		rules:     make(RuleMap),
		resources: make(ResourceMap),
	}

	rules := []RuleSpec{}
	err := yaml2.Unmarshal(GetReferenceRules(), &rules)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		control.rules[rule.GroupVersionKind] = rule
	}

	return control, nil
}

// Map is
func (rc *ResourceController) Map() ResourceMap {
	return rc.resources
}

// LoadResourcesFormBytes is
func (rc *ResourceController) LoadResourcesFormBytes(in []byte) ([]*resource.Resource, error) {
	var err error
	var result []*resource.Resource

	defaultconfig.GetDefaultFieldSpecs()

	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(in), 1024)

	for err == nil || isEmptyYamlError(err) {
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

type refreshFunc func(interface{}) (interface{}, error)

func (rc *ResourceController) refreshMap(resourceMap map[string]interface{}, keys []string, create bool, fn refreshFunc) (err error) {
	if len(keys) == 0 {
		return nil
	}

	key := keys[0]
	if _, found := resourceMap[key]; !found {
		if !create {
			return nil
		}

		resourceMap[key] = map[string]interface{}{}
	}

	keys = keys[1:]
	if len(keys) == 0 {
		resourceMap[key], err = fn(resourceMap[key])

		return err
	}

	switch nextVal := resourceMap[key].(type) {
	case nil:
		log.Printf("nil value at `%s` ignored in mutation attempt", strings.Join(keys, "."))
		return nil
	case []interface{}:
		for i := range nextVal {

			if nextMap, ok := nextVal[i].(map[string]interface{}); !ok {
				return fmt.Errorf("%#v is expected to be %T", nextVal[i], nextMap)
			} else if err := rc.refreshMap(nextMap, keys, create, fn); err != nil {
				return err
			}
		}
		return nil
	case map[string]interface{}:
		return rc.refreshMap(nextVal, keys, create, fn)
	default:
		return fmt.Errorf("%#v is not expected to be a primitive type", nextVal)
	}
}

func (rc *ResourceController) refreshFields(res *resource.Resource, fields []FieldSpec, fn func(resource.GroupVersionKind, FieldPath, interface{}) (interface{}, error)) error {
	for _, field := range fields {
		for _, path := range field.Paths {
			err := rc.refreshMap(res.Object, path.Slice(), field.Create, func(in interface{}) (interface{}, error) {
				return fn(field.GroupVersionKind, path, in)
			})

			if nil != err {
				return err
			}
		}
	}

	return nil
}

// FindResource is
func (rc *ResourceController) FindResource(id resource.ID) *resource.Resource {
	for id, res := range rc.resources {
		if id.Equals(id) {
			return res
		}

	}

	return nil
}

// RefreshReferences is
func (rc *ResourceController) RefreshReferences() (err error) {
	for key, _ := range rc.rules {
		fmt.Println(key)
	}

	for _, res := range rc.resources {
		fmt.Print("res.ID() => ")
		fmt.Println(res.ID())
		references, exists := rc.rules[res.ID()]
		if !exists {
			continue
		}

		err = rc.refreshFields(res, references.MatedataName, func(gvk resource.GroupVersionKind, path FieldPath, in interface{}) (interface{}, error) {
			return res.GetName(), nil
		})
		if nil != err {
			return err
		}

		err = rc.refreshFields(res, references.MatedataLabels, func(gvk resource.GroupVersionKind, path FieldPath, in interface{}) (interface{}, error) {
			return res.GetLabels(), nil
		})
		if nil != err {
			return err
		}

		err = rc.refreshFields(res, references.MatedataAnnotations, func(gvk resource.GroupVersionKind, path FieldPath, in interface{}) (interface{}, error) {
			return res.GetAnnotations(), nil
		})
		if nil != err {
			return err
		}
	}

	return nil
}
