package resource

import (
	"strings"
)

const (
	escapedForwardSlash  = "\\/"
	tempSlashReplacement = "???"
)

// Path is
type Path string

func (p Path) String() string {
	return string(p)
}

// Slice is
func (p Path) Slice() []string {
	path := p.String()

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
