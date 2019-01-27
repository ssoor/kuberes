package target

import (
	"io/ioutil"

	"github.com/ssoor/kuberes/pkg/gvk"
	"github.com/ssoor/kuberes/pkg/reference"
	"github.com/ssoor/kuberes/pkg/yaml"
)

// ReferenceRule is
type ReferenceRule struct {
	gvk.GVK `json:",inline,omitempty" yaml:",inline,omitempty"`

	MatedataName        []reference.FieldSpec `json:"matedata.name,omitempty" yaml:"matedata.name,omitempty"`
	MatedataLabels      []reference.FieldSpec `json:"matedata.labels,omitempty" yaml:"matedata.labels,omitempty"`
	MatedataAnnotations []reference.FieldSpec `json:"metadata.annotations,omitempty" yaml:"metadata.annotations,omitempty"`
}

// ReferenceRuleMap is a map from resource ID to Resource.
type ReferenceRuleMap map[gvk.GVK]ReferenceRule

// Load is
func (r ReferenceRuleMap) Load(path string) error {
	body, err := ioutil.ReadFile(path)
	if nil != err {
		return err
	}
	out := ReferenceRule{}
	decoder := yaml.NewFormatErrorDecodeFormBytes(body, path)

	for err == nil {
		if err = decoder.Decode(&out); nil != err {
			continue
		}

		err = nil
		r[out.GVK] = out
	}

	return nil
}
