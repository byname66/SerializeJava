package structures

import (
	"fmt"
	"main/common"
)

type ClassDescInfo struct {
	ClassDescFlags  byte
	Fields          *Fields
	ClassAnnotation *Annotation
	SuperClassDesc  *ClassDesc
}

func NewClassDescInfo(classDescFlags byte, fields *Fields, Annotation *Annotation, superClassDesc *ClassDesc) *ClassDescInfo {
	return &ClassDescInfo{
		ClassDescFlags:  classDescFlags,
		Fields:          fields,
		ClassAnnotation: Annotation,
		SuperClassDesc:  superClassDesc,
	}
}

func ParseClassDescInfo(parser *StructuresParser) (*ClassDescInfo, error) {

	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, err
	}
	classDescFlags := signByte
	fields, err := ParseFields(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseClassDescInfo:\n %v", err)
	}
	Annotation, err := ParseAnnotation(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseClassDescInfo:\n %v", err)
	}
	superClassDesc, err := ParseClassDesc(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseClassDescInfo:\n %v", err)
	}
	return NewClassDescInfo(classDescFlags, fields, Annotation, superClassDesc), nil
}

// check if the input value is contained in classDescFlags
func (classDescInfo ClassDescInfo) HasFlag(flag byte) bool {

	return (classDescInfo.ClassDescFlags & flag) == flag
}

func (classDescInfo *ClassDescInfo) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteByte(classDescInfo.ClassDescFlags)
	err := classDescInfo.Fields.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteClassDescInfo:\n %v", err)
	}
	err = classDescInfo.ClassAnnotation.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteClassDescInfo:\n %v", err)
	}
	err = classDescInfo.SuperClassDesc.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteClassDescInfo:\n %v", err)
	}
	return nil
}

func (CDI *ClassDescInfo) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	switch CDI.ClassDescFlags {
	case 0x02:
		str = "SC_SERIALIZABLE"
	case 0x03:
		str = "SC_SERIALIZABLE|SC_WRITE_METHOD"
	case 0x04:
		str = "SC_EXTERNALIZABLE"
	case 0x0C:
		str = "SC_EXTERNALIZABLE|SC_BLOCK_DATA"
	default:
		str = fmt.Sprintf("%v", CDI.ClassDescFlags)
		// default:
		// 	return "", fmt.Errorf("in ClassDescInfo#ToString:No matched Flag")
	}
	result += sb.Buildf("- classDescFlags  ", []interface{}{CDI.ClassDescFlags, "  -  ", str})
	str, err = CDI.Fields.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in ClassDescInfo#ToString:\n%v", err)
	}
	result += str
	str, err = CDI.ClassAnnotation.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in ClassDescInfo#ToString:\n%v", err)
	}
	result += str
	str, err = CDI.SuperClassDesc.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in ClassDescInfo#ToString:\n%v", err)
	}
	result += str
	return result, nil
}
