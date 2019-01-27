package target

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/kubernetes-sigs/kustomize/pkg/transformers/config/defaultconfig"
	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// ResourceMap is a map from resource ID to Resource.
type ResourceMap map[resource.UniqueID]*resource.Resource

// Bytes encodes a ResMap to YAML; encoded objects separated by `---`.
func (rm ResourceMap) Bytes() ([]byte, error) {
	var ids []resource.UniqueID
	for id := range rm {
		ids = append(ids, id)
	}
	// sort.Sort(IdSlice(ids))

	buf := bytes.NewBuffer([]byte{})

	for _, id := range ids {
		out, err := yaml.Marshal(rm[id].Map())
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

// Add is
func (rm ResourceMap) Add(override bool, id resource.UniqueID, res *resource.Resource) error {
	if !override {
		if _, found := rm[id]; found {
			return fmt.Errorf("id '%q' already used", id)
		}
	}

	rm[id] = res

	return nil
}

// Range is
func (rm ResourceMap) Range(fn func(resource.UniqueID, *resource.Resource) error) error {
	for id, res := range rm {
		if err := fn(id, res); nil != err {
			return err
		}
	}

	return nil
}

// Merge is
func (rm ResourceMap) Merge(override bool, resources ...ResourceMap) error {
	for _, m := range resources {
		if m == nil {
			continue
		}

		for id, res := range m {
			if err := rm.Add(override, id, res); nil != err {
				return err
			}
		}
	}

	return nil
}

// MergeFormResource is
func (rm ResourceMap) MergeFormResource(override bool, resources ...*resource.Resource) error {
	for _, res := range resources {
		if err := rm.Add(override, res.ID(), res); nil != err {
			return err
		}
	}

	return nil
}

// MergeFormPath is
func (rm ResourceMap) MergeFormPath(override bool, path string) error {
	body, err := ioutil.ReadFile(path)
	if nil != err {
		return err
	}

	resources, err := rm.loadResources(yaml.NewFormatErrorDecodeFromBytes(body, path))
	if nil != err {
		return err
	}

	if err := rm.MergeFormResource(override, resources...); nil != err {
		return err
	}

	return nil
}

func (rm ResourceMap) loadResources(decoder yaml.Decoder) ([]*resource.Resource, error) {
	var err error
	var result []*resource.Resource

	defaultconfig.GetDefaultFieldSpecs()

	for err == nil {
		var out unstructured.Unstructured

		if err = decoder.Decode(&out); err != nil {
			continue
		}

		res := &resource.Resource{Unstructured: out}

		if err := res.Validate(); err != nil {
			return nil, err
		}

		result = append(result, res)
	}

	if err != io.EOF {
		return nil, err
	}

	return result, nil
}
