package merge

import "k8s.io/apimachinery/pkg/util/strategicpatch"

type strategicpatchMerge struct {
}

// Merge applies a strategic merge patch. The patch and the original document
// must be json encoded content. A patch can be created from an original and a modified document
// by calling CreateStrategicMergePatch.
func (m *strategicpatchMerge) Merge(original, patch []byte, dataStruct interface{}) ([]byte, error) {
	return strategicpatch.StrategicMergePatch(original, patch, dataStruct)
}

// MergeMap applies a strategic merge patch. The original and patch documents
// must be JSONMap. A patch canbe  created from an original and modified document by
// calling CreateTwoWayMergeMapPatch.
// Warning: the original and patch JSONMap objects are mutated by this function and should not be reused.
func (m *strategicpatchMerge) MergeMap(original, patch JSONMap, dataStruct interface{}) (JSONMap, error) {
	merged, err := strategicpatch.StrategicMergeMapPatch(strategicpatch.JSONMap(original), strategicpatch.JSONMap(patch), dataStruct)

	return JSONMap(merged), err
}
