package target

import (
	"github.com/ssoor/kuberes/pkg/resource"
)

// FieldSpec is
type FieldSpec struct {
	resource.GVK `json:",inline,omitempty" yaml:",inline,omitempty"`

	Force bool            `json:"create,omitempty" yaml:"create,omitempty"`
	Paths []resource.Path `json:"paths,omitempty" yaml:"paths,omitempty"`
}

// FieldSpecCallback is
type FieldSpecCallback func(resource.Path, interface{}) (interface{}, error)

// Refresh is
func (fs FieldSpec) Refresh(res *resource.Resource, fn FieldSpecCallback) error {
	for _, path := range fs.Paths {
		err := res.ScanPath(path, fs.Force, func(data interface{}) (interface{}, error) {
			return fn(path, data)
		})

		if nil != err {
			return err
		}
	}

	return nil
}
