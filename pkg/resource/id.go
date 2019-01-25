package resource

import (
	"strings"
)

const (
	noNamespace = "~X"
	noPrefix    = "~P"
	noName      = "~N"
	noSuffix    = "~S"
	separator   = "|"
)

// ID conflates GroupVersionKind with a textual name to uniquely identify a kubernetes resource (object).
type ID interface {
	String() string
	Equals(ID) bool
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

// GroupVersionKind unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GroupVersionKind struct {
	Kind    string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Group   string `json:"group,omitempty" yaml:"group,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

func (fs GroupVersionKind) String() string {
	return generationIDString(fs.Group, fs.Version, fs.Kind)
}

// Equals returns true if the Gvk's have equal fields.
func (fs GroupVersionKind) Equals(o ID) bool {
	return fs.String() == o.String()
}
