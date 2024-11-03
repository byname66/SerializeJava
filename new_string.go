package structures

import (
	"fmt"
)

type NewString struct {
	tc_string byte
	newHandle uint32
	length    int16
	value     string
}

func NewNewString(handle uint32, length int16, value string) *NewString {
	return &NewString{
		tc_string: TC_STRING,
		newHandle: handle,
		length:    length,
		value:     value,
	}
}

func ParseNewString(parser *StructuresParser) (*NewString, error) {
	signByte, err := parser.byteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("In ParseClassName1:\n %w", err)
	}
	if signByte != TC_STRING {
		return nil, fmt.Errorf("In ParseClassName1:No TC_STRING")
	}
	length, value, err := parser.ReadName()
	if err != nil {
		return nil, fmt.Errorf("In ParseClassName1:\n %w", err)
	}
	handle := parser.AddHandle()

	return NewNewString(handle, length, value), nil
}
