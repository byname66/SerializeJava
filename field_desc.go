package structures

import "fmt"

type FieldDesc struct {
	primitiveDesc *PrimitiveDesc
	objectDesc    *ObjectDesc
}

func ParseFieldDesc(parser *StructuresParser) (*FieldDesc, error) {
	peekByte, err := parser.byteReader.PeekByte()
	if err != nil {
		return nil, fmt.Errorf("In ParseFieldDesc:\n +%w", err)
	}
	fieldDesc := new(FieldDesc)
	switch string(peekByte) {
	case "B", "C", "D", "F", "I", "J", "S", "Z":
		fieldDesc.primitiveDesc, err = ParsePrimitiveDesc(parser)
	case "[", "L":
		fieldDesc.objectDesc, err = ParseObjectDesc(parser)
	default:
		return nil, fmt.Errorf("In ParseFieldDesc:No prim_typecode or obj_typecode")
	}
	if err != nil {
		return nil, fmt.Errorf("In ParseFieldDesc:\n +%w", err)
	}
	return fieldDesc, nil
}
