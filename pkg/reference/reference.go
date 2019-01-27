package reference

import (
	"github.com/ssoor/kuberes/pkg/resource"
)

// Reference is
type Reference struct {
	resource.GVK `json:",inline,omitempty" yaml:",inline,omitempty"`

	MatedataName        []FieldSpec `json:"matedata.name,omitempty" yaml:"matedata.name,omitempty"`
	MatedataLabels      []FieldSpec `json:"matedata.labels,omitempty" yaml:"matedata.labels,omitempty"`
	MatedataAnnotations []FieldSpec `json:"metadata.annotations,omitempty" yaml:"metadata.annotations,omitempty"`
}

// RefreshCallback is
type RefreshCallback func(FieldSpec, FieldPath, interface{}) (interface{}, error)

// RefreshMatedataName is
func (i Reference) RefreshMatedataName(res *resource.Resource, fn RefreshCallback) error {
	return i.refreshFields(res, i.MatedataName, fn)
}

// RefreshMatedataLabels is
func (i Reference) RefreshMatedataLabels(res *resource.Resource, fn RefreshCallback) error {
	return i.refreshFields(res, i.MatedataLabels, fn)
}

// RefreshMatedataAnnotations is
func (i Reference) RefreshMatedataAnnotations(res *resource.Resource, fn RefreshCallback) error {
	return i.refreshFields(res, i.MatedataAnnotations, fn)
}

func (i Reference) refreshFields(res *resource.Resource, fields []FieldSpec, fn RefreshCallback) error {
	for _, field := range fields {
		if err := field.Refresh(res.Map(), fn); nil != err {
			return err
		}
	}

	return nil
}
