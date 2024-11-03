package structures

import (
	"fmt"
)

type ObjectDesc struct {
	obj_typecode string
	fieldName    *FieldName
	className1   *NewString
}

func NewObjectDesc(obj_typecode string, fieldName *FieldName, className1 *NewString) *ObjectDesc {
	return &ObjectDesc{
		obj_typecode: obj_typecode,
		fieldName:    fieldName,
		className1:   className1,
	}
}

func ParseObjectDesc(parser *StructuresParser) (*ObjectDesc, error) {
	typeByte, err := parser.byteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("In ParsePrimitiveDesc:\n %w", err)
	}
	var obj_typecode string
	switch string(typeByte) {
	case "[":
		obj_typecode = "Array"
	case "L":
		obj_typecode = "Object"
	default:
		return nil, fmt.Errorf("In ParsePrimitiveDesc:No obj_typecode")
	}
	fieldName, err := ParseFieldName(parser)
	if err != nil {
		return nil, fmt.Errorf("In ParseObjectDesc:\n %w", err)
	}
	className1, err := ParseNewString(parser)
	if err != nil {
		return nil, fmt.Errorf("In ParseObjectDesc:\n %w", err)
	}
	return NewObjectDesc(obj_typecode, fieldName, className1), nil
}
