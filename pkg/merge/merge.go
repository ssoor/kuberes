package merge

// Merger is
type Merger interface {
	// Merge applies a strategic merge patch. The patch and the original document
	// must be json encoded content. A patch can be created from an original and a modified document
	// by calling CreateStrategicMergePatch.
	Merge(original []byte, dataStruct interface{}) ([]byte, error)
	// MergeMap applies a strategic merge patch. The original and patch documents
	// must be JSONMap. A patch can be created from an original and modified document by
	// calling CreateTwoWayMergeMapPatch.
	// Warning: the original and patch JSONMap objects are mutated by this function and should not be reused.
	MergeMap(original JSONMap, dataStruct interface{}) (JSONMap, error)
}

// NewJSONMerger is
func NewJSONMerger(patch JSONMap) Merger {
	return &jsonpatchMerge{patch: patch}
}

// NewRFC6902Merger is
func NewRFC6902Merger(patch []byte) Merger {
	return &rfc6902Merge{patch: patch}
}

// NewStrategicMerge is
func NewStrategicMerge(patch JSONMap) Merger {
	return &strategicpatchMerge{patch: patch}
}
