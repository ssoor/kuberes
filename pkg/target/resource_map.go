package target

import (
	"bytes"

	"github.com/ghodss/yaml"
	"github.com/ssoor/kuberes/pkg/resource"
)

// ResourceMap is a map from resource ID to Resource.
type ResourceMap map[resource.UniqueID]*resource.Resource

// Yaml encodes a ResMap to YAML; encoded objects separated by `---`.
func (rm ResourceMap) Yaml() ([]byte, error) {
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
