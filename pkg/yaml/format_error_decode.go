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

// NewFormatErrorDecode is
func NewFormatErrorDecode(de Decoder, path string) Decoder {
	return &formatErrorDecode{de: de, path: path}
}

// NewFormatErrorDecodeFormBytes is
func NewFormatErrorDecodeFormBytes(body []byte, path string) Decoder {
	return &formatErrorDecode{de: NewYAMLOrJSONDecoderFormBytes(body), path: path}
}
