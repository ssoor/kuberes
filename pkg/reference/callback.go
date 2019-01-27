package reference

import "github.com/ssoor/kuberes/pkg/resource"

// RefreshSpec is
type RefreshSpec struct {
	resource.GVK

	Name     string
	Path     FieldPath
	Resource *resource.Resource
}

// RefreshCallback is
type RefreshCallback func(RefreshSpec, interface{}) (interface{}, error)
