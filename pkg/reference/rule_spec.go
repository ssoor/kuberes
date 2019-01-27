package reference

import "github.com/ssoor/kuberes/pkg/resource"

// RuleSpec is
type RuleSpec struct {
	resource.GVKID `json:",inline,omitempty" yaml:",inline,omitempty"`

	MatedataName        []FieldSpec `json:"matedata.name,omitempty" yaml:"matedata.name,omitempty"`
	MatedataLabels      []FieldSpec `json:"matedata.labels,omitempty" yaml:"matedata.labels,omitempty"`
	MatedataAnnotations []FieldSpec `json:"metadata.annotations,omitempty" yaml:"metadata.annotations,omitempty"`
}

// Refresh is
func (rs RuleSpec) Refresh(res *resource.Resource, fn func(RuleSpec, FieldSpec, FieldPath, interface{}) (interface{}, error)) (err error) {
	err = rs.refreshFields(res, rs.MatedataName, func(fs FieldSpec, fp FieldPath, data interface{}) (interface{}, error) {
		return fn(rs, fs, fp, data)
	})
	if nil != err {
		return err
	}

	err = rs.refreshFields(res, rs.MatedataLabels, func(fs FieldSpec, fp FieldPath, data interface{}) (interface{}, error) {
		return fn(rs, fs, fp, data)
	})
	if nil != err {
		return err
	}

	err = rs.refreshFields(res, rs.MatedataAnnotations, func(fs FieldSpec, fp FieldPath, data interface{}) (interface{}, error) {
		return fn(rs, fs, fp, data)
	})
	if nil != err {
		return err
	}

	return nil
}

func (rs RuleSpec) refreshFields(res *resource.Resource, fields []FieldSpec, fn func(FieldSpec, FieldPath, interface{}) (interface{}, error)) error {
	for _, field := range fields {
		if err := field.Refresh(res, fn); nil != err {
			return err
		}
	}

	return nil
}
