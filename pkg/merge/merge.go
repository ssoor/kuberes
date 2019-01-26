package merge

// JSONMap is a representations of JSON object encoded as map[string]interface{}
// where the children can be either map[string]interface{}, []interface{} or
// primitive type).
// Operating on JSONMap representation is much faster as it doesn't require any
// json marshaling and/or unmarshaling operations.
type JSONMap map[string]interface{}

// Merger is
type Merger interface {
	// Merge applies a strategic merge patch. The patch and the original document
	// must be json encoded content. A patch can be created from an original and a modified document
	// by calling CreateStrategicMergePatch.
	Merge(original, patch []byte, dataStruct interface{}) ([]byte, error)
	// MergeMap applies a strategic merge patch. The original and patch documents
	// must be JSONMap. A patch can be created from an original and modified document by
	// calling CreateTwoWayMergeMapPatch.
	// Warning: the original and patch JSONMap objects are mutated by this function and should not be reused.
	MergeMap(original, patch JSONMap, dataStruct interface{}) (JSONMap, error)
}
