package target

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"sigs.k8s.io/kustomize/pkg/constants"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/types"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	"github.com/ssoor/kuberes/pkg/log"
)

// Target is
type Target struct {
	ldr           ifc.Loader
	kustomization *types.Kustomization
}

// NewTarget is
func NewTarget(ldr ifc.Loader) (*Target, error) {
	content, err := loadKustFile(ldr)
	if err != nil {
		return nil, err
	}

	var k types.Kustomization
	if err := unmarshal(content, &k); err != nil {
		return nil, err
	}

	k.DealWithDeprecatedFields()
	msgs, errs := k.EnforceFields()
	if len(errs) > 0 {
		return nil, fmt.Errorf(strings.Join(errs, "\n"))
	}
	if len(msgs) > 0 {
		log.Printf(strings.Join(msgs, "\n"))
	}

	return &Target{ldr: ldr, kustomization: &k}, nil
}

func loadKustFile(ldr ifc.Loader) ([]byte, error) {
	for _, kf := range []string{
		constants.KustomizationFileName,
		constants.SecondaryKustomizationFileName} {
		content, err := ldr.Load(kf)
		if err == nil {
			return content, nil
		}
		if !strings.Contains(err.Error(), "no such file or directory") {
			return nil, err
		}
	}
	return nil, fmt.Errorf("no kustomization.yaml file under %s", ldr.Root())
}

func unmarshal(y []byte, o interface{}) error {
	j, err := yaml.YAMLToJSON(y)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(j))
	dec.DisallowUnknownFields()
	return dec.Decode(o)
}

// YamlFormatError represents error with yaml file name where json/yaml format error happens.
type YamlFormatError struct {
	Path     string
	ErrorMsg string
}

func (e YamlFormatError) Error() string {
	return fmt.Sprintf("YAML file [%s] encounters a format error.\n%s\n", e.Path, e.ErrorMsg)
}

func isYAMLSyntaxError(e error) bool {
	return strings.Contains(e.Error(), "error converting YAML to JSON") || strings.Contains(e.Error(), "error unmarshaling JSON")
}

// ErrorHandler handles YamlFormatError
func ErrorHandler(e error, path string) error {
	if isYAMLSyntaxError(e) {
		return YamlFormatError{
			Path:     path,
			ErrorMsg: e.Error(),
		}
	}
	return e
}

// LoadResources is
func (t *Target) LoadResources() (*ResourceController, error) {
	controller, err := NewResourceController()
	if nil != err {
		return nil, err
	}

	for _, path := range t.kustomization.Resources {
		content, err := t.ldr.Load(path)
		if err != nil {
			return nil, errors.Wrap(err, "Load from path "+path+" failed")
		}

		res, err := controller.LoadResourcesFormBytes(content)
		if err != nil {
			return nil, ErrorHandler(err, path)
		}

		if err := controller.MergeResources(false, res...); nil != err {
			return nil, err
		}
	}

	return controller, nil
}
