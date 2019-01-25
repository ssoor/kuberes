package resource

// Matedata allows manipulation of k8s objects
// that do not have Golang structs.
type Matedata interface {
	GetName() string
	SetName(string)
	GetNamespace() string
	SetNamespace(string)

	GetLabels() map[string]string
	SetLabels(map[string]string)
	GetAnnotations() map[string]string
	SetAnnotations(map[string]string)
}
