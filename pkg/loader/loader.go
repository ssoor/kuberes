package loader

import (
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/loader"

	"github.com/ssoor/kuberes/pkg/yaml"
)

// Loader interface exposes methods to read bytes.
type Loader interface {
	// Root returns the root location for this Loader.
	Root() string
	// New returns Loader located at newRoot.
	Sub(path string) Loader
	// LoadBytes returns the bytes read from the location or an error.
	LoadBytes(path string) ([]byte, error)
	// LoadYamlDecoder returns the yaml.Decoder read from the location or an error.
	LoadYamlDecoder(path string) (yaml.Decoder, error)
	// Cleanup cleans the loader
	Cleanup() error
}

// NewLoader returns a Loader.
func NewLoader(root string, fSys fs.FileSystem) (Loader, error) {
	ldr, err := loader.NewLoader(root, fSys)
	if nil != err {
		return nil, err
	}

	return newFileLoader(ldr), nil
}
