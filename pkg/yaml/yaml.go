package yaml

import (
	"bytes"
	"io"

	yaml2 "github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	defaultBufferSize = 1024
)

// Decoder is a convenience interface for Decode.
type Decoder interface {
	Decode(into interface{}) error
}

// Marshal the object into JSON then converts JSON to YAML and returns the
// YAML.
func Marshal(o interface{}) ([]byte, error) {
	return yaml2.Marshal(o)
}

// ToJSONFromBytes converts YAML to JSON. Since JSON is a subset of YAML,
// passing JSON through this method should be a no-op.
//
// Things YAML can do that are not supported by JSON:
// * In YAML you can have binary and null keys in your maps. These are invalid
//   in JSON. (int and float keys are converted to strings.)
// * Binary data in YAML with the !!binary tag is not supported. If you want to
//   use binary data with this library, encode the data as base64 as usual but do
//   not use the !!binary tag in your YAML. This will ensure the original base64
//   encoded data makes it all the way through to the JSON.
//
// For strict decoding of YAML, use YAMLToJSONStrict.
func ToJSONFromBytes(y []byte) ([]byte, error) {
	return yaml2.YAMLToJSON(y)
}

// NewYAMLOrJSONDecoderFromBytes returns a decoder that will process YAML documents
// or JSON documents from the given reader as a stream. bufferSize determines
// how far into the stream the decoder will look to figure out whether this
// is a JSON stream (has whitespace followed by an open brace).
func NewYAMLOrJSONDecoderFromBytes(body []byte) Decoder {
	return yaml.NewYAMLOrJSONDecoder(bytes.NewReader(body), defaultBufferSize)
}

// NewYAMLOrJSONDecoderFromReader returns a decoder that will process YAML documents
// or JSON documents from the given reader as a stream. bufferSize determines
// how far into the stream the decoder will look to figure out whether this
// is a JSON stream (has whitespace followed by an open brace).
func NewYAMLOrJSONDecoderFromReader(r io.Reader, bufferSize int) Decoder {
	return yaml.NewYAMLOrJSONDecoder(r, bufferSize)
}
