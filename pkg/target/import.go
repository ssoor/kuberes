package target

import "github.com/ssoor/kuberes/pkg/resourcemap"

// Import is
type Import struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Attach string `json:"attach,omitempty" yaml:"attach,omitempty"`
}

// Make is
func (i Import) Make() (resourcemap.ResourceMap, error) {
	t, err := NewTarget()
	if nil != err {
		return nil, err
	}

	if err := t.Load(i.Attach); err != nil {
		return nil, err
	}
	if err := t.Make(); err != nil {
		return nil, err
	}

	return t.ResourceMap(), nil
}
