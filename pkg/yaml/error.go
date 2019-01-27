package yaml

import "fmt"

// Error represents error with yaml file name where json/yaml format error happens.
type Error struct {
	Err  error
	Path string
}

func (e Error) Error() string {
	return fmt.Sprintf("YAML file [%s] encounters a format error.\n%s\n", e.Path, e.Err.Error())
}
