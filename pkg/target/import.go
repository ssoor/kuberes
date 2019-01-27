package target

import "github.com/ssoor/kuberes/pkg/loader"

// Import is
type Import struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Attach string `json:"attach,omitempty" yaml:"attach,omitempty"`
}

// Make is
func (i Import) Make(loader loader.Loader) (ResourceMap, error) {
	t, err := NewTarget(loader.Sub(i.Attach))
	if nil != err {
		return nil, err
	}

	if err := t.Load(); err != nil {
		return nil, err
	}

	t.Name = i.Name

	return t.Make()
}
