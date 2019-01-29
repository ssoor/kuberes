package target

import (
	"github.com/ssoor/kuberes/pkg/loader"
	"github.com/ssoor/kuberes/pkg/resource"
)

// Import is
type Import struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Attach string `json:"attach,omitempty" yaml:"attach,omitempty"`
}

// Make is
func (i Import) Make(loader loader.Loader) (map[resource.UniqueID]*resource.Resource, error) {
	t, err := NewMaker(loader.Sub(i.Attach), i.Name)
	if nil != err {
		return nil, err
	}

	return t.Make()
}
