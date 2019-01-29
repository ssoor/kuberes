package target

// Import is
type Import struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Attach string `json:"attach,omitempty" yaml:"attach,omitempty"`
}

// Patchs is
type Patchs struct {
	RFC6902   []RFC6902Patch `json:"rfc6902,omitempty" yaml:"rfc6902,omitempty"`
	Strategic []string       `json:"strategic,omitempty" yaml:"strategic,omitempty"`
}

// Matedata is
type Matedata struct {
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// Target is
type Target struct {
	Name    string   `json:"name,omitempty" yaml:"name,omitempty"`
	Imports []Import `json:"imports,omitempty" yaml:"imports,omitempty"`

	// Patchs to add to all objects.
	Patchs Patchs `json:"patchs,omitempty" yaml:"patchs,omitempty"`

	// Matedata to add to all objects.
	Matedata Matedata `json:"matedata,omitempty" yaml:"matedata,omitempty"`

	// Resources specifies relative paths to files holding YAML representations
	// of kubernetes API objects. URLs and globs not supported.
	Resources []string `json:"resources,omitempty" yaml:"resources,omitempty"`
}
