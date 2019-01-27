package reference

import "github.com/ssoor/kuberes/pkg/resource"

// RefreshSpec is
type RefreshSpec struct {
	resource.GVK `json:",inline,omitempty" yaml:",inline,omitempty"`

	Name      string    `json:"name,omitempty" yaml:"name,omitempty"`
	FieldPath FieldPath `json:"path,omitempty" yaml:"path,omitempty"`
}

// RefreshCallback is
type RefreshCallback func(RefreshSpec, interface{}) (interface{}, error)
