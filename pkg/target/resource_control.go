package target

import (
	"bytes"
	"fmt"

	"github.com/ssoor/kuberes/pkg/loader"
	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
)

// ResourceController is
type ResourceController interface {
	Bytes() ([]byte, error)

	Map() map[resource.UniqueID]*resource.Resource

	Get(id resource.UniqueID) *resource.Resource
	Range(fn func(resource.UniqueID, *resource.Resource) error) error
	Add(id resource.UniqueID, res *resource.Resource, override bool) error
}

// resourceControl is a map from resource ID to Resource.
type resourceControl struct {
	loader      loader.Loader
	resourceMap map[resource.UniqueID]*resource.Resource
}

// NewResourceControl is
func NewResourceControl(loader loader.Loader, paths []string) (ResourceController, error) {
	rc := &resourceControl{
		loader:      loader,
		resourceMap: make(map[resource.UniqueID]*resource.Resource),
	}

	for _, path := range paths {
		decoder, err := loader.LoadYamlDecoder(path)
		if nil != err {
			return nil, err
		}

		resources, err := resource.NewResourceFromDecoder(decoder)
		if nil != err {
			return nil, err
		}

		for _, res := range resources {

			if err := rc.Add(res.ID(), res, false); nil != err {
				return nil, err
			}
		}
	}

	return rc, nil
}

// Map is
func (rc resourceControl) Map() map[resource.UniqueID]*resource.Resource {
	return rc.resourceMap
}

// Get is
func (rc resourceControl) Get(id resource.UniqueID) *resource.Resource {
	return rc.resourceMap[id]
}

// Add is
func (rc resourceControl) Add(id resource.UniqueID, res *resource.Resource, override bool) error {
	if !override {
		if _, found := rc.resourceMap[id]; found {
			return fmt.Errorf("id '%q' already used", id)
		}
	}

	rc.resourceMap[id] = res

	return nil
}

// Range is
func (rc resourceControl) Range(fn func(resource.UniqueID, *resource.Resource) error) error {
	for id, res := range rc.resourceMap {
		if err := fn(id, res); nil != err {
			return err
		}
	}

	return nil
}

// Merge is
func (rc resourceControl) Merge(override bool, resources ...*resource.Resource) error {
	for _, res := range resources {
		if err := rc.Add(res.ID(), res, override); nil != err {
			return err
		}
	}

	return nil
}

// Bytes encodes a ResMap to YAML; encoded objects separated by `---`.
func (rc resourceControl) Bytes() ([]byte, error) {
	var ids []resource.UniqueID
	for id := range rc.resourceMap {
		ids = append(ids, id)
	}
	// sort.Sort(IdSlice(ids))

	buf := bytes.NewBuffer([]byte{})

	for _, id := range ids {
		out, err := yaml.Marshal(rc.Get(id).Map())
		if err != nil {
			return nil, err
		}

		if _, err = buf.Write(out); err != nil {
			return nil, err
		}

		if _, err = buf.Write([]byte("---\n")); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
