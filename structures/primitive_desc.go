package structures

import (
	"fmt"
	"main/common"
)

type PrimitiveDesc struct {
	Prim_typecode string
	FieldName     *UTF
}

func NewPrimitiveDesc(prim_typecode string, fieldName *UTF) *PrimitiveDesc {
	return &PrimitiveDesc{
		Prim_typecode: prim_typecode,
		FieldName:     fieldName,
	}
}
func ParsePrimitiveDesc(parser *StructuresParser) (*PrimitiveDesc, error) {
	typeByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParsePrimitiveDesc:\n %v", err)
	}
	fieldName, err := ParseUtf(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParsePrimitiveDesc:\n %v", err)
	}
	return NewPrimitiveDesc(string(typeByte), fieldName), nil
}

func (pd *PrimitiveDesc) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteNByte([]byte(pd.Prim_typecode))
	err := pd.FieldName.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in ParsePrimitiveDesc:\n %v", err)
	}
	return nil
}

func (pd *PrimitiveDesc) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @PrimitiveDesc")
	indent += IndentSpaceCount
	result += sb.Buildf("- typeCode  ", []interface{}{pd.Prim_typecode, "  -  ", []byte(pd.Prim_typecode)})
	str, err = pd.FieldName.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in PrimitiveDesc:\n%v", err)
	}
	result += str
	return result, nil
}
