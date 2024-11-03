package structures

import "fmt"

type FieldName struct {
	length int16
	value  string
}

func NewFieldName(length int16, value string) *FieldName {
	return &FieldName{
		length: length,
		value:  value,
	}
}

func ParseFieldName(parser *StructuresParser) (*FieldName, error) {
	length, value, err := parser.ReadName()
	if err != nil {
		return nil, fmt.Errorf("In ParseFieldName:\n %w", err)
	}
	return NewFieldName(length, value), nil
}
