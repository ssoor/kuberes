package reference

import (
	"fmt"
	"strings"

	"github.com/ssoor/kuberes/pkg/log"
)

const (
	escapedForwardSlash  = "\\/"
	tempSlashReplacement = "???"
)

// FieldPath is
type FieldPath string

func (fp FieldPath) String() string {
	return string(fp)
}

// Slice is
func (fp FieldPath) Slice() []string {
	path := fp.String()

	if !strings.Contains(path, escapedForwardSlash) {
		return strings.Split(path, "/")
	}
	s := strings.Replace(path, escapedForwardSlash, tempSlashReplacement, -1)
	paths := strings.Split(s, "/")
	var result []string
	for _, path := range paths {
		result = append(result, strings.Replace(path, tempSlashReplacement, "/", -1))
	}
	return result
}

// Refresh is
func (fp FieldPath) Refresh(body map[string]interface{}, create bool, fn func(FieldPath, interface{}) (interface{}, error)) error {
	return fp.refreshMap(body, fp.Slice(), create, fn)
}

func (fp FieldPath) refreshMap(resourceMap map[string]interface{}, keys []string, create bool, fn func(FieldPath, interface{}) (interface{}, error)) (err error) {
	if len(keys) == 0 {
		return nil
	}

	key := keys[0]
	if _, found := resourceMap[key]; !found {
		if !create {
			return nil
		}

		resourceMap[key] = map[string]interface{}{}
	}

	keys = keys[1:]
	if len(keys) == 0 {
		resourceMap[key], err = fn(fp, resourceMap[key])

		return err
	}

	switch nextVal := resourceMap[key].(type) {
	case nil:
		log.Printf("nil value at `%s` ignored in mutation attempt", strings.Join(keys, "."))
		return nil
	case []interface{}:
		for i := range nextVal {

			if nextMap, ok := nextVal[i].(map[string]interface{}); !ok {
				return fmt.Errorf("%#v is expected to be %T", nextVal[i], nextMap)
			} else if err := fp.refreshMap(nextMap, keys, create, fn); err != nil {
				return err
			}
		}
		return nil
	case map[string]interface{}:
		return fp.refreshMap(nextVal, keys, create, fn)
	default:
		return fmt.Errorf("%#v is not expected to be a primitive type", nextVal)
	}
}
