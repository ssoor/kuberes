package merge

import "encoding/json"

// JSONMap is a representations of JSON object encoded as map[string]interface{}
// where the children can be either map[string]interface{}, []interface{} or
// primitive type).
// Operating on JSONMap representation is much faster as it doesn't require any
// json marshaling and/or unmarshaling operations.
type JSONMap map[string]interface{}

// NewJSONMap is
func NewJSONMap(j []byte) (map[string]interface{}, error) {
	if j == nil {
		j = []byte("{}")
	}

	m := map[string]interface{}{}
	err := json.Unmarshal(j, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Marshal is
func (jm JSONMap) Marshal() ([]byte, error) {
	return json.Marshal(jm)
}
