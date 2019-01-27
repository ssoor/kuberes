package reference

import (
	"io"

	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
)

// Reference is
type Reference struct {
	resource.GVK `json:",inline,omitempty" yaml:",inline,omitempty"`

	// MatedataName        []FieldSpec `json:"matedata.name,omitempty" yaml:"matedata.name,omitempty"`
	// MatedataLabels      []FieldSpec `json:"matedata.labels,omitempty" yaml:"matedata.labels,omitempty"`
	// MatedataAnnotations []FieldSpec `json:"metadata.annotations,omitempty" yaml:"metadata.annotations,omitempty"`

	FieldSpecs map[string][]FieldSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// NewReferenceFromDecoder is
func NewReferenceFromDecoder(decoder yaml.Decoder) ([]*Reference, error) {
	var err error
	var result []*Reference

	for err == nil {
		out := &Reference{}

		if err = decoder.Decode(out); err != nil {
			continue
		}

		result = append(result, out)
	}

	if err != io.EOF {
		return nil, err
	}

	return result, nil
}

// Refresh is
func (i Reference) Refresh(res *resource.Resource, fn RefreshCallback) error {
	for name, specs := range i.FieldSpecs {
		for _, field := range specs {
			if err := field.Refresh(res.Map(), func(rs RefreshSpec, data interface{}) (interface{}, error) {
				rs.Name = name
				rs.Resource = res

				return fn(rs, data)
			}); nil != err {
				return err
			}
		}
	}

	return nil
}
