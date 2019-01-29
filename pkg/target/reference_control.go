package target

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ssoor/kuberes/pkg/loader"
	"github.com/ssoor/kuberes/pkg/log"
	"github.com/ssoor/kuberes/pkg/resource"
	"github.com/ssoor/kuberes/pkg/yaml"
)

// ReferenceRule is
type ReferenceRule struct {
	resource.GVK `json:",inline,omitempty" yaml:",inline,omitempty"`

	FieldSpecs map[string][]FieldSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// ReferenceControl is a map from resource ID to Resource.
type ReferenceControl struct {
	handleMap    map[string]RefreshHandle
	referenceMap map[resource.GVK]*ReferenceRule
}

// RefreshHandle is
type RefreshHandle func(FieldSpec, resource.Path, interface{}) (interface{}, error)

// NewReferenceControl is
func NewReferenceControl(l loader.Loader, path string) (*ReferenceControl, error) {
	rc := &ReferenceControl{
		handleMap:    make(map[string]RefreshHandle),
		referenceMap: make(map[resource.GVK]*ReferenceRule),
	}

	decoder, err := l.LoadYamlDecoder(path)
	if nil != err {
		body, err := ioutil.ReadFile(path)
		if nil != err {
			return nil, err
		}

		decoder = yaml.NewFormatErrorDecodeFromBytes(body, path)
	}

	var items []*ReferenceRule

	err = nil
	for err == nil {
		out := &ReferenceRule{}

		if err = decoder.Decode(out); err != nil {
			continue
		}

		items = append(items, out)
	}

	if err != io.EOF {
		return nil, err
	}

	for _, ref := range items {
		log.Debug("REFERENCE ADD:", ref.GVK)
		rc.referenceMap[ref.GVK] = ref
	}

	return rc, nil
}

// AddRefreshHandle is
func (rc *ReferenceControl) AddRefreshHandle(name string, handle RefreshHandle) {
	rc.handleMap[name] = handle
}

// GetRules is
func (rc ReferenceControl) getRule(key resource.GVK) *ReferenceRule {
	item, exists := rc.referenceMap[key]
	if !exists {
		return nil // continue
	}

	return item
}

// Refresh is
func (rc ReferenceControl) Refresh(res *resource.Resource) error {
	rule := rc.getRule(res.GVK())
	if nil == rule {
		log.Debug("REFERENCE NOT FOUND:", res.GVK())
		return nil
	}

	for name, fileds := range rule.FieldSpecs {
		log.Debug("REFERENCE:", name, res.GVK())

		for _, field := range fileds {
			err := field.Refresh(res, func(path resource.Path, in interface{}) (out interface{}, err error) {
				handle := rc.handleMap[name]

				if out, err = handle(field, path, in); nil != err {
					return nil, err
				}

				fmt.Printf("[%s] %s:%s(%v) => %v\n", res.ID(), name, path, in, out)

				return out, nil
			})

			if nil != err {
				return err
			}
		}
	}

	return nil
}
