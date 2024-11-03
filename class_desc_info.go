package structures

import (
	"fmt"
)

type ClassDescInfo struct {
	classDescFlags  byte
	Fields          *Fields
	classAnnotation *ClassAnnotation
	superClassDesc  *ClassDesc
}

func NewClassDescInfo(fields *Fields, classAnnotation *ClassAnnotation, superClassDesc *ClassDesc) *ClassDescInfo {
	return &ClassDescInfo{
		Fields:          fields,
		classAnnotation: classAnnotation,
		superClassDesc:  superClassDesc,
	}
}

func ParseClassDescInfo(parser *StructuresParser) (*ClassDescInfo, error) {

	signByte, err := parser.byteReader.ReadByte()
	if err != nil {
		return nil, err
	}
	if signByte != FLAG_CLASS_DESC {
		return nil, fmt.Errorf("In ParseClassDescInfo: No ClassDescFlags")
	}
	fields, err := ParseFields(parser)
	if err != nil {
		return nil, fmt.Errorf("In ParseClassDescInfo:\n %w", err)
	}
	classAnnotation, err := ParseClassAnnotation(parser)
	if err != nil {
		return nil, fmt.Errorf("In ParseClassDescInfo:\n %w", err)
	}
	superClassDesc, err := ParseClassDesc(parser)
	if err != nil {
		return nil, fmt.Errorf("In ParseClassDescInfo:\n %w", err)
	}
	return NewClassDescInfo(fields, classAnnotation, superClassDesc), nil
}
