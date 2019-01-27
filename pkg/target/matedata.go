package target

import "github.com/ssoor/kuberes/pkg/resource"

// Matedata is
type Matedata struct {
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// Make is
func (m Matedata) Make(res *resource.Resource) error {
	res.SetNamespace(m.Namespace)
	res.SetLabels(m.MergeMap(res.GetLabels(), m.Labels))
	res.SetAnnotations(m.MergeMap(res.GetAnnotations(), m.Annotations))

	return nil
}

// MergeMap is
func (m Matedata) MergeMap(src map[string]string, merge map[string]string) map[string]string {
	if nil == merge {
		return src
	}

	if nil == src {
		return merge
	}

	for key, val := range merge {
		src[key] = val
	}

	return src
}
