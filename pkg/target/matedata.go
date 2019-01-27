package target

import "github.com/ssoor/kuberes/pkg/resource"

// Matedata is
type Matedata struct {
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// MakeResource is
func (m Matedata) MakeResource(res *resource.Resource) {
	res.SetNamespace(m.Namespace)
	res.SetLabels(m.Labels)
	res.SetAnnotations(m.Annotations)
}
