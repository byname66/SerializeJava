package structures

import (
	"fmt"
	"main/common"
)

type ObjectDesc struct {
	OBJ_typecode string
	FieldName    *UTF
	ClassName1   *NewString
}

func NewObjectDesc(obj_typecode string, fieldName *UTF, className1 *NewString) *ObjectDesc {
	return &ObjectDesc{
		OBJ_typecode: obj_typecode,
		FieldName:    fieldName,
		ClassName1:   className1,
	}
}

func ParseObjectDesc(parser *StructuresParser) (*ObjectDesc, error) {
	typeByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParsePrimitiveDesc:\n %v", err)
	}
	fieldName, err := ParseUtf(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseObjectDesc:\n %v", err)
	}
	className1, err := ParseNewString(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseObjectDesc:\n %v", err)
	}
	return NewObjectDesc(string(typeByte), fieldName, className1), nil
}

func (od *ObjectDesc) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteNByte([]byte(od.OBJ_typecode))
	err := od.FieldName.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteObjectDesc:\n %v", err)
	}
	err = od.ClassName1.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteObjectDesc:\n %v", err)
	}
	return nil
}

func (od *ObjectDesc) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @ObjectDesc")
	indent += IndentSpaceCount
	str, err = od.FieldName.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in ObjectDesc#ToString:\n%v", err)
	}
	result += str
	str, err = od.ClassName1.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in ObjectDesc#ToString:\n%v", err)
	}
	result += str
	return result, nil
}
