package resource

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Header allows manipulation of k8s objects
// that do not have Golang structs.
type Header interface {
	GetKind() string
	SetKind(string)
	GetAPIVersion() string
	SetAPIVersion(string)
	GroupVersionKind() schema.GroupVersionKind
	SetGroupVersionKind(schema.GroupVersionKind)
}
