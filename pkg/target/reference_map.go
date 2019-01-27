package target

import (
	"io/ioutil"

	"github.com/ssoor/kuberes/pkg/gvk"
	"github.com/ssoor/kuberes/pkg/reference"
	"github.com/ssoor/kuberes/pkg/yaml"
)

// ReferenceMap is a map from resource ID to Resource.
type ReferenceMap map[gvk.GVK]reference.Reference

// Load is
func (r ReferenceMap) Load(path string) error {
	body, err := ioutil.ReadFile(path)
	if nil != err {
		return err
	}

	out := reference.Reference{}
	decoder := yaml.NewFormatErrorDecodeFromBytes(body, path)

	for err == nil {
		if err = decoder.Decode(&out); nil != err {
			continue
		}

		err = nil
		r[out.GVK] = out
	}

	return nil
}

// FindByGVK is
func (r ReferenceMap) FindByGVK(key gvk.GVK) *reference.Reference {
	item, exists := r[key]
	if !exists {
		return nil // continue
	}

	return &item
}
