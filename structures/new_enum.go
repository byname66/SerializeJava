package structures

import (
	"fmt"
	"main/common"
)

type NewEnum struct {
	TC_ENUM          byte
	ClassDesc        *ClassDesc
	NewHandle        uint32
	EnumConstantName *EnumConstantName
}

func NewNewEnum(ClassDesc *ClassDesc, newHandle uint32, enumConstantName *EnumConstantName) *NewEnum {
	return &NewEnum{
		TC_ENUM:          TC_ENUM,
		ClassDesc:        ClassDesc,
		NewHandle:        newHandle,
		EnumConstantName: enumConstantName,
	}
}

func ParseNewEnum(parser *StructuresParser) (*NewEnum, error) {
	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseNewEnum:\n %v", err)
	}
	if signByte != TC_ENUM {
		return nil, fmt.Errorf("in ParseNewEnum: NO TC_ENUM")
	}
	classDesc, err := ParseClassDesc(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseNewEnum:\n %v", err)
	}
	newHandle := parser.AddHandle()
	EnumConstantName, err := ParseEnumConstantName(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseNewEnum:\n %v", err)
	}
	newEnum := NewNewEnum(classDesc, newHandle, EnumConstantName)
	parser.AddReferenced(newEnum)
	return newEnum, nil
}

func (ne *NewEnum) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteByte(TC_ENUM)
	err := ne.ClassDesc.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteNewEnum:\n %v", err)
	}
	err = ne.EnumConstantName.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteNewEnum:\n %v", err)
	}
	return nil
}

func (ne *NewEnum) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @NewEnum")
	indent += IndentSpaceCount
	result += sb.Buildf("- TC_ENUM  ", []interface{}{TC_ENUM})
	str, err = ne.ClassDesc.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in NewEnum#ToString\n%v", err)
	}
	result += str
	result += sb.Buildf("- newHandle  ", []interface{}{ne.NewHandle})
	str, err = ne.EnumConstantName.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in NewEnum#ToString\n%v", err)
	}
	result += str
	return result, nil
}

func (ne *NewEnum) GetNewHandle() uint32 {
	return ne.NewHandle
}

type EnumConstantName struct {
	StringObject *Object
}

func ParseEnumConstantName(parser *StructuresParser) (*EnumConstantName, error) {
	stringObject, err := ParseObject(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseEnumConstantName:\n %v", err)
	}
	EnumConstantName := new(EnumConstantName)
	EnumConstantName.StringObject = stringObject
	return EnumConstantName, nil
}

func (ecn *EnumConstantName) ToByte(parser *StructuresParser) error {
	err := ecn.StringObject.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteEnumConstantName:\n %v", err)
	}
	return nil
}

func (ecn *EnumConstantName) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @EnumConstantName")
	indent += IndentSpaceCount
	str, err = ecn.StringObject.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in EnumConstantName#ToString\n%v", err)
	}
	result += str
	return result, nil
}
