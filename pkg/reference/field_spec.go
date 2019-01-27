package reference

import "github.com/ssoor/kuberes/pkg/resource"

// FieldSpec is
type FieldSpec struct {
	resource.GVK `json:",inline,omitempty" yaml:",inline,omitempty"`

	Create bool        `json:"create,omitempty" yaml:"create,omitempty"`
	Paths  []FieldPath `json:"paths,omitempty" yaml:"paths,omitempty"`
}

// Refresh is
func (fs FieldSpec) Refresh(body map[string]interface{}, fn RefreshCallback) error {
	for _, path := range fs.Paths {
		err := path.Refresh(body, fs.Create, func(fp FieldPath, data interface{}) (interface{}, error) {
			return fn(RefreshSpec{GVK: fs.GVK, Path: path}, data)
		})

		if nil != err {
			return err
		}
	}

	return nil
}
