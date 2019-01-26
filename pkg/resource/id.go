package resource

import (
	"strings"
)

// ID conflates GroupVersionKind with a textual name to uniquely identify a kubernetes resource (object).
type ID interface {
	String() string
	Equals(ID) bool
}

// GVKID unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GVKID struct {
	Group   string `json:"group,omitempty" yaml:"group,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Kind    string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

func (fs GVKID) String() string {
	return generationIDString(fs.Group, fs.Version, fs.Kind)
}

// Equals returns true if the Gvk's have equal fields.
func (fs GVKID) Equals(o ID) bool {
	return fs.String() == o.String()
}

func generationIDString(group, version, kind string) string {
	// Values that are brief but meaningful in logs.
	const (
		noGroup   = "~G"
		noVersion = "~V"
		noKind    = "~K"
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
