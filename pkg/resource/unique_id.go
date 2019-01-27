package resource

import (
	"strings"
)

// UniqueID unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type UniqueID struct {
	GVK

	// original name of the resource before transformation.
	Name string
}

// NewUniqueID is
func NewUniqueID(name string, gvk GVK) UniqueID {
	return UniqueID{Name: name, GVK: gvk}
}

func (uid UniqueID) String() string {
	return generationIDString(uid.Name, uid.GVK)
}

// Equals returns true if the Gvk's have equal fields.
func (uid UniqueID) Equals(o UniqueID) bool {
	return uid.String() == o.String()
}

func generationIDString(name string, gvk GVK) string {
	// Values that are brief but meaningful in logs.
	const (
		noName    = "~N"
		separator = "_"
	)

	if name == "" {
		name = noName
	}

	return strings.Join([]string{noName, gvk.String()}, separator)
}
