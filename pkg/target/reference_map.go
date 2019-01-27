package target

import (
	"github.com/ssoor/kuberes/pkg/reference"
	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
)

// ReferenceMap is a map from resource ID to Resource.
type ReferenceMap map[resource.GVK]*reference.Reference

// Load is
func (r ReferenceMap) Load(decoder yaml.Decoder) (err error) {
	references, err := reference.NewReferenceFromDecoder(decoder)
	if nil != err {
		return err
	}

	for _, ref := range references {
		r[ref.GVK] = ref
	}

	return nil
}

// FindByGVK is
func (r ReferenceMap) FindByGVK(key resource.GVK) *reference.Reference {
	item, exists := r[key]
	if !exists {
		return nil // continue
	}

	return item
}
