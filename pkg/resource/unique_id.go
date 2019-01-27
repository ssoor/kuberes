package resource

import (
	"strings"

	"github.com/ssoor/kuberes/pkg/gvk"
)

// UniqueID unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type UniqueID struct {
	gvk.GVK

	// original name of the resource before transformation.
	Name string
	// namespace the resource belongs to
	// an untransformed resource has no namespace, fully transformed resource has the namespace from
	// the top most overlay
	Namespace string
}

// NewUniqueID is
func NewUniqueID(name, namespace string, gvk gvk.GVK) UniqueID {
	return UniqueID{Name: name, Namespace: namespace, GVK: gvk}
}

func (uid UniqueID) String() string {
	return generationIDString(uid.Name, uid.Namespace, uid.GVK)
}

// Equals returns true if the Gvk's have equal fields.
func (uid UniqueID) Equals(o UniqueID) bool {
	return uid.String() == o.String()
}

func generationIDString(name, namespace string, gvk gvk.GVK) string {
	// Values that are brief but meaningful in logs.
	const (
		noName      = "~N"
		noNamespace = "~NS"
		separator   = "_"
	)

	if name == "" {
		name = noName
	}

	if namespace == "" {
		namespace = noNamespace
	}

	return strings.Join([]string{noName, noNamespace, gvk.String()}, separator)
}
