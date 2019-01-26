package target

import (
	"strings"

	"github.com/ssoor/kuberes/pkg/resource"
)

const (
	escapedForwardSlash  = "\\/"
	tempSlashReplacement = "???"
)

// FieldPath is
type FieldPath string

func (rs FieldPath) String() string {
	return string(rs)
}

// Slice is
func (rs FieldPath) Slice() []string {
	path := rs.String()

	if !strings.Contains(path, escapedForwardSlash) {
		return strings.Split(path, "/")
	}
	s := strings.Replace(path, escapedForwardSlash, tempSlashReplacement, -1)
	paths := strings.Split(s, "/")
	var result []string
	for _, path := range paths {
		result = append(result, strings.Replace(path, tempSlashReplacement, "/", -1))
	}
	return result
}

// FieldSpec is
type FieldSpec struct {
	resource.GVKID `json:",inline,omitempty" yaml:",inline,omitempty"`

	Create bool        `json:"create,omitempty" yaml:"create,omitempty"`
	Paths  []FieldPath `json:"paths,omitempty" yaml:"paths,omitempty"`
}

// RuleSpec is
type RuleSpec struct {
	resource.GVKID `json:",inline,omitempty" yaml:",inline,omitempty"`

	MatedataName        []FieldSpec `json:"matedata.name,omitempty" yaml:"matedata.name,omitempty"`
	MatedataLabels      []FieldSpec `json:"matedata.labels,omitempty" yaml:"matedata.labels,omitempty"`
	MatedataAnnotations []FieldSpec `json:"metadata.annotations,omitempty" yaml:"metadata.annotations,omitempty"`
}
