package structures

import "fmt"

type Fields struct {
	Count      int16
	FieldDescs []*FieldDesc
}

func NewFields(count int16, fieldDescs []*FieldDesc) *Fields {
	return &Fields{
		Count:      count,
		FieldDescs: fieldDescs,
	}
}

func ParseFields(parser *StructuresParser) (*Fields, error) {
	count, err := parser.byteReader.ReadInt16()
	if err != nil {
		return nil, fmt.Errorf("In ParseFields:\n %w", err)
	}
	var fieldDescs []*FieldDesc
	for i := 0; i < int(count); i++ {
		fieldsDesc, err := ParseFieldDesc(parser)
		if err != nil {
			return nil, fmt.Errorf("In ParseFields:\n %w", err)
		}
		fieldDescs = append(fieldDescs, fieldsDesc)
	}
	return NewFields(count, fieldDescs), nil
}
