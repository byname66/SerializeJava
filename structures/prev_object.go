package structures

import (
	"fmt"
	"main/common"
)

type PrevObject struct {
	TC_REFERENCE     byte
	Handler          uint32
	ReferencedObject ReferencedObject
}

func NewPrevObject(handle uint32, obj ReferencedObject) *PrevObject {
	return &PrevObject{
		TC_REFERENCE:     TC_REFERENCE,
		Handler:          handle,
		ReferencedObject: obj,
	}
}

func ParsePrevObject(parser *StructuresParser) (*PrevObject, error) {
	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParsePrevObject:\n %v", err)
	}
	if signByte != TC_REFERENCE {
		return nil, fmt.Errorf("in ParsePrevObject: NO TC_REFERENCE")
	}
	handle, err := parser.ByteReader.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("in ParsePrevObject:\n %v", err)
	}
	handle0 := uint32(handle)
	obj, err := parser.GetReferenced(handle0)
	if err != nil {
		return nil, fmt.Errorf("in ParsePrevObject:\n %v", err)
	}
	return NewPrevObject(handle0, obj), nil
}

func (po *PrevObject) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteByte(TC_REFERENCE)
	parser.ByteReader.WriteNumber(po.Handler)
	return nil
}

func (po *PrevObject) ToString(indent int) (string, error) {
	var (
		result string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @PrevObject")
	indent += IndentSpaceCount
	result += sb.Buildf("- TC_REFERENCE  ", []interface{}{TC_REFERENCE})
	byteArray, err := common.ConvertNumberToBytes(po.Handler)
	if err != nil {
		return "", fmt.Errorf("in PrevObject#ToString:\n%v", err)
	}
	result += sb.Buildf("- handle", []interface{}{po.Handler, "  -  ", byteArray})
	return result, nil
}
