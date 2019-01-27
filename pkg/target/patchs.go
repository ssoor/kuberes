package target

import "github.com/ssoor/kuberes/pkg/resource"

// Patchs is
type Patchs struct {
	Base      []string                 `json:"base,omitempty" yaml:"base,omitempty"`
	RFC6902   []map[string]interface{} `json:"rfc6902,omitempty" yaml:"rfc6902,omitempty"`
	Strategic []string                 `json:"strategic,omitempty" yaml:"strategic,omitempty"`
}

// MakeResource is
func (m Patchs) MakeResource(res *resource.Resource) {
}
