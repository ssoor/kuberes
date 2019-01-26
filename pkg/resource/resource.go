package resource

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// GenerationBehavior specifies generation behavior of configmaps, secrets and maybe other resources.
type GenerationBehavior int

const (
	// BehaviorUnspecified is an Unspecified behavior; typically treated as a Create.
	BehaviorUnspecified GenerationBehavior = iota
	// BehaviorCreate makes a new resource.
	BehaviorCreate
	// BehaviorReplace replaces a resource.
	BehaviorReplace
	// BehaviorMerge attempts to merge a new resource with an existing resource.
	BehaviorMerge
)

// String converts a GenerationBehavior to a string.
func (b GenerationBehavior) String() string {
	switch b {
	case BehaviorReplace:
		return "replace"
	case BehaviorMerge:
		return "merge"
	case BehaviorCreate:
		return "create"
	default:
		return "unspecified"
	}
}

// NewGenerationBehavior converts a string to a GenerationBehavior.
func NewGenerationBehavior(s string) GenerationBehavior {
	switch s {
	case "replace":
		return BehaviorReplace
	case "merge":
		return BehaviorMerge
	case "create":
		return BehaviorCreate
	default:
		return BehaviorUnspecified
	}
}

// Resource is map representation of a Kubernetes API resource object
// paired with a GenerationBehavior.
type Resource struct {
	unstructured.Unstructured

	b GenerationBehavior
}

// String returns resource as JSON.
func (r *Resource) String() string {
	bs, err := r.MarshalJSON()
	if err != nil {
		return "<" + err.Error() + ">"
	}
	return r.b.String() + ":" + strings.TrimSpace(string(bs))
}

// ID returns the ID for the resource.
func (r *Resource) ID() ID {
	gvk := r.GroupVersionKind()

	return GVKID{gvk.Group, gvk.Version, gvk.Kind}
}

// Map returns the Map for the resource.
func (r *Resource) Map() map[string]interface{} {
	return r.Object
}

// Header returns the header for the resource.
func (r *Resource) Header() Header {
	return r
}

// Matedata returns the matedata for the resource.
func (r *Resource) Matedata() Matedata {
	return r
}

// Behavior returns the behavior for the resource.
func (r *Resource) Behavior() GenerationBehavior {
	return r.b
}

// SetBehavior changes the resource to the new behavior
func (r *Resource) SetBehavior(b GenerationBehavior) *Resource {
	r.b = b
	return r
}

// Validate validates that u has kind and name
// except for kind `List`, which doesn't require a name
func (r *Resource) Validate() error {
	kind := r.GetKind()
	if kind == "" {
		return fmt.Errorf("missing kind in object %v", r.Unstructured)
	} else if kind == "List" {
		return nil
	}
	if r.GetName() == "" {
		return fmt.Errorf("missing metadata.name in object %v", r.Unstructured)
	}

	return nil
}

// IsGenerated checks if the resource is generated from a generator
func (r *Resource) IsGenerated() bool {
	return r.b != BehaviorUnspecified
}

// Merge performs merge with other resource.
func (r *Resource) Merge(other *Resource) {
	r.Replace(other)
	mergeConfigmap(r.Map(), other.Map(), r.Map())
}

// Replace performs replace with other resource.
func (r *Resource) Replace(other *Resource) {
	r.SetLabels(mergeStringMaps(other.GetLabels(), r.GetLabels()))
	r.SetAnnotations(
		mergeStringMaps(other.GetAnnotations(), r.GetAnnotations()))
	r.SetName(other.GetName())
}

// TODO: Add BinaryData once we sync to new k8s.io/api
func mergeConfigmap(
	mergedTo map[string]interface{},
	maps ...map[string]interface{}) {
	mergedMap := map[string]interface{}{}
	for _, m := range maps {
		datamap, ok := m["data"].(map[string]interface{})
		if ok {
			for key, value := range datamap {
				mergedMap[key] = value
			}
		}
	}
	mergedTo["data"] = mergedMap
}

func mergeStringMaps(maps ...map[string]string) map[string]string {
	result := map[string]string{}
	for _, m := range maps {
		for key, value := range m {
			result[key] = value
		}
	}
	return result
}
