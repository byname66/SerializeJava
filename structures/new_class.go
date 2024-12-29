package structures

import (
	"fmt"
	"main/common"
)

type NewClass struct {
	TC_class  byte
	ClassDesc *ClassDesc
	NewHandle uint32
}

func NewNewClass(classDesc *ClassDesc, newHandle uint32) *NewClass {
	return &NewClass{
		TC_class:  TC_CLASS,
		ClassDesc: classDesc,
		NewHandle: newHandle,
	}
}

func ParseNewClass(parser *StructuresParser) (*NewClass, error) {
	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseNewClass:\n %v", err)
	}
	if signByte != TC_CLASS {
		return nil, fmt.Errorf("in ParseNewClass:No TC_CLASS")
	}
	classDesc, err := ParseClassDesc(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseNewClass:\n %v", err)
	}
	newHandle := parser.AddHandle()
	newClass := NewNewClass(classDesc, newHandle)
	parser.AddReferenced(newClass)
	return newClass, nil
}

func (nc *NewClass) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteByte(TC_CLASS)
	err := nc.ClassDesc.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteNewClass:\n %v", err)
	}
	return nil
}

func (nc *NewClass) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @NewClass")
	indent += IndentSpaceCount
	result += sb.Buildf("- TC_CLASS  ", []interface{}{TC_CLASS})
	str, err = nc.ClassDesc.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in NewClass#ToString\n%v", err)
	}
	result += str
	result += sb.Buildf("- newHandle  ", []interface{}{nc.NewHandle})
	return result, nil
}

func (nc *NewClass) GetNewHandle() uint32 {
	return nc.NewHandle
}
