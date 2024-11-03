package structures

import "fmt"

type NewClass struct {
	tc_class  byte
	classDesc *ClassDesc
	newHandle uint32
}

func NewNewClass(classDesc *ClassDesc, newHandle uint32) *NewClass {
	return &NewClass{
		tc_class:  TC_CLASS,
		classDesc: classDesc,
		newHandle: newHandle,
	}
}

func ParseNewClass(parser *StructuresParser) (*NewClass, error) {
	signByte, err := parser.byteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("In ParseNewClass:\n %w", err)
	}
	if signByte != TC_CLASS {
		return nil, fmt.Errorf("In ParseNewClass:No TC_CLASS")
	}
	classDesc, err := ParseClassDesc(parser)
	if err != nil {
		return nil, fmt.Errorf("In ParseNewClass:\n %w", err)
	}
	newHandle := parser.AddHandle()
	return NewNewClass(classDesc, newHandle), nil
}
