package merge

import (
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"
)

type jsonpatchMerge struct {
	patch JSONMap
}

// Merge applies a strategic merge patch. The patch and the original document
// must be json encoded content. A patch can be created from an original and a modified document
// by calling CreateStrategicMergePatch.
func (m *jsonpatchMerge) Merge(original []byte, dataStruct interface{}) ([]byte, error) {
	// Use JSON merge patch to handle types w/o schema
	patchBytes, err := json.Marshal(m.patch)
	if err != nil {
		return nil, err
	}

	return jsonpatch.MergePatch(original, patchBytes)
}

// MergeMap applies a strategic merge patch. The original and patch documents
// must be JSONMap. A patch can be created from an original and modified document by
// calling CreateTwoWayMergeMapPatch.
// Warning: the original and patch JSONMap objects are mutated by this function and should not be reused.
func (m *jsonpatchMerge) MergeMap(original JSONMap, dataStruct interface{}) (merged JSONMap, err error) {
	// Use JSON merge patch to handle types w/o schema
	originalBytes, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}

	mergedBytes, err := m.Merge(originalBytes, dataStruct)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(mergedBytes, &merged)
	if err != nil {
		return nil, err
	}

	return merged, nil
}
