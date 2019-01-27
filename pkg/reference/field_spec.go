package reference

import "github.com/ssoor/kuberes/pkg/resource"

// FieldSpec is
type FieldSpec struct {
	resource.GVKID `json:",inline,omitempty" yaml:",inline,omitempty"`

	Create bool        `json:"create,omitempty" yaml:"create,omitempty"`
	Paths  []FieldPath `json:"paths,omitempty" yaml:"paths,omitempty"`
}

// Refresh is
func (fs FieldSpec) Refresh(res *resource.Resource, fn func(FieldSpec, FieldPath, interface{}) (interface{}, error)) error {
	for _, path := range fs.Paths {
		err := path.Refresh(res, fs.Create, func(fp FieldPath, data interface{}) (interface{}, error) {
			return fn(fs, fp, data)
		})

		if nil != err {
			return err
		}
	}

	return nil
}