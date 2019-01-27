package merge

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

type strategicpatchMerge struct {
	patch JSONMap
}

// Merge applies a strategic merge patch. The patch and the original document
// must be json encoded content. A patch can be created from an original and a modified document
// by calling CreateStrategicMergePatch.
func (m *strategicpatchMerge) Merge(original []byte, dataStruct interface{}) ([]byte, error) {
	originalMap, err := NewJSONMap(original)
	if nil != err {
		return nil, err
	}

	result, err := m.MergeMap(originalMap, dataStruct)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

// MergeMap applies a strategic merge patch. The original and patch documents
// must be JSONMap. A patch canbe  created from an original and modified document by
// calling CreateTwoWayMergeMapPatch.
// Warning: the original and patch JSONMap objects are mutated by this function and should not be reused.
func (m *strategicpatchMerge) MergeMap(original JSONMap, dataStruct interface{}) (JSONMap, error) {
	merged, err := strategicpatch.StrategicMergeMapPatch(strategicpatch.JSONMap(original), strategicpatch.JSONMap(m.patch), dataStruct)

	return JSONMap(merged), err
}
