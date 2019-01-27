package yaml

import "io"

type formatErrorDecode struct {
	de   Decoder
	path string
}

func (d *formatErrorDecode) Decode(into interface{}) error {
	err := d.de.Decode(into)

	if nil != err {
		if io.EOF == err {
			return err
		}

		return Error{Err: err, Path: d.path}
	}

	return nil
}

// NewFromatErrorDecode is
func NewFromatErrorDecode(de Decoder, path string) Decoder {
	return &formatErrorDecode{de: de, path: path}
}

// NewFormatErrorDecodeFromBytes is
func NewFormatErrorDecodeFromBytes(body []byte, path string) Decoder {
	return &formatErrorDecode{de: NewYAMLOrJSONDecoderFromBytes(body), path: path}
}
