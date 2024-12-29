package structures

import (
	"fmt"
	"main/common"
)

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
	count, err := parser.ByteReader.ReadInt16()
	if err != nil {
		return nil, fmt.Errorf("in ParseFields:\n %v", err)
	}
	var fieldDescs []*FieldDesc
	for i := 0; i < int(count); i++ {
		fieldsDesc, err := ParseFieldDesc(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseFields:\n %v", err)
		}
		fieldDescs = append(fieldDescs, fieldsDesc)
	}
	return NewFields(count, fieldDescs), nil
}

func (f *Fields) ToByte(parser *StructuresParser) error {
	err := parser.ByteReader.WriteNumber(f.Count)
	if err != nil {
		return fmt.Errorf("in WriteFields:\n %v", err)
	}
	for i := 0; i < int(f.Count); i++ {
		err = f.FieldDescs[i].ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteFields:\n %v", err)
		}
	}
	return nil
}

func (f *Fields) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @Fields")
	indent += IndentSpaceCount
	byteArray, err := common.ConvertNumberToBytes(f.Count)
	if err != nil {
		return "", fmt.Errorf("in Fields#ToString:\n%v", err)
	}
	result += sb.Buildf("- count  ", []interface{}{f.Count, "  -  ", byteArray})
	indent += IndexToArraySpaceCount
	for i := 0; i < len(f.FieldDescs); i++ {
		result += sb.Buildf("Index  [", []interface{}{i, "]"})
		str, err = f.FieldDescs[i].ToString(indent + IndentSpaceCount)
		if err != nil {
			return "", fmt.Errorf("in Fields#ToString:\n%v", err)
		}
		result += str
	}
	return result, nil
}
