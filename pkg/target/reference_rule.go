package target

import (
	"io/ioutil"

	"github.com/ssoor/kuberes/pkg/reference"
	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
)

// ReferenceRule is
type ReferenceRule struct {
	resource.GVKID `json:",inline,omitempty" yaml:",inline,omitempty"`

	MatedataName        []reference.FieldSpec `json:"matedata.name,omitempty" yaml:"matedata.name,omitempty"`
	MatedataLabels      []reference.FieldSpec `json:"matedata.labels,omitempty" yaml:"matedata.labels,omitempty"`
	MatedataAnnotations []reference.FieldSpec `json:"metadata.annotations,omitempty" yaml:"metadata.annotations,omitempty"`
}

// ReferenceRuleMap is a map from resource ID to Resource.
type ReferenceRuleMap map[resource.GVKID]ReferenceRule

// LoadReferenceRuleMapFormFile is
func LoadReferenceRuleMapFormFile(filename string) (ReferenceRuleMap, error) {
	body, err := ioutil.ReadFile(filename)
	if nil != err {
		return nil, err
	}

	return LoadReferenceRuleMapFormBytes(body, filename)
}

// LoadReferenceRuleMapFormBytes is
func LoadReferenceRuleMapFormBytes(body []byte, path string) (ruleMap ReferenceRuleMap, err error) {
	ruleMap = make(ReferenceRuleMap)

	out := ReferenceRule{}
	decoder := yaml.NewFormatErrorDecodeFormBytes(body, path)

	for err == nil {
		if err = decoder.Decode(&out); nil != err {
			continue
		}

		err = nil
		ruleMap[out.GVKID] = out
	}

	return ruleMap, nil
}
