package structures

import (
	"bytes"
	"fmt"
	"io"
)

type Stream struct {
	Magic    []byte
	Version  []byte
	Contents *[]Content
}

func NewStream(Contents *[]Content) *Stream {
	return &Stream{
		Magic:    STREAM_MAGIC,
		Version:  STREAM_VERSION,
		Contents: Contents,
	}
}

func NewEmptyStream() *Stream {
	return &Stream{
		Magic:    STREAM_MAGIC,
		Version:  STREAM_VERSION,
		Contents: nil,
	}
}

// first to check the Magic and Version bytes
func CheckMagicAndVersion(parser *StructuresParser) error {
	tbytes, error := parser.byteReader.ReadNByte(2)
	if error != nil {
		return fmt.Errorf("error reading Magic bytes: %w\n", error)
	}
	if !bytes.Equal(tbytes, STREAM_MAGIC) {
		return fmt.Errorf("Magic bytes do not match\n")
	}
	tbytes, error = parser.byteReader.ReadNByte(2)
	if error != nil {
		return fmt.Errorf("error reading Version bytes: %w\n", error)
	}
	if !bytes.Equal(tbytes, STREAM_VERSION) {
		return fmt.Errorf("Version bytes do not match\n")
	}
	return nil
}

func ParseStream(parser *StructuresParser) (*Stream, error) {
	err := CheckMagicAndVersion(parser)
	if err != nil {
		return nil, err
	}
	contents := new([]Content)
	for {

		content, err := ParseContent(parser)
		if err != nil {
			if err == io.EOF {
				return NewStream(contents), nil
			}
			return nil, err
		}

		*contents = append(*contents, *content)
	}

}
