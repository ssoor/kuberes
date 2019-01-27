package loader

import (
	"sigs.k8s.io/kustomize/pkg/ifc"

	"github.com/ssoor/kuberes/pkg/yaml"
)

type fileLoader struct {
	ifc.Loader
}

func newFileLoader(ldr ifc.Loader) *fileLoader {
	return &fileLoader{Loader: ldr}
}

func (l fileLoader) Sub(path string) Loader {
	ldr, err := l.Loader.New(path)
	if nil != err {
		panic("TODO")
	}

	return newFileLoader(ldr)
}

func (l fileLoader) LoadBytes(path string) ([]byte, error) {
	return l.Loader.Load(path)
}

func (l fileLoader) LoadYamlDecoder(path string) (yaml.Decoder, error) {
	body, err := l.LoadBytes(path)
	if nil != err {
		return nil, err
	}

	return yaml.NewFormatErrorDecodeFromBytes(body, path), nil
}
