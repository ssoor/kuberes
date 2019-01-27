package resource

import (
	"strings"
)

// GVK unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GVK struct {
	Group   string `json:"group,omitempty" yaml:"group,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Kind    string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

func (fs GVK) String() string {
	return generationGVKIDString(fs.Group, fs.Version, fs.Kind)
}

// Equals returns true if the Gvk's have equal fields.
func (fs GVK) Equals(o GVK) bool {
	return fs.String() == o.String()
}

func generationGVKIDString(group, version, kind string) string {
	// Values that are brief but meaningful in logs.
	const (
		noGroup   = "{}"
		noVersion = "{}"
		noKind    = "{}"
		separator = "_"
	)

	if group == "" {
		group = noGroup
	}
	if version == "" {
		version = noVersion
	}
	if kind == "" {
		kind = noKind
	}

	return strings.Join([]string{group, version, kind}, separator)
}
