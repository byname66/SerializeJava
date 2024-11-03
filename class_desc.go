package structures

import "fmt"

type ClassDesc struct {
	newClassDesc  *NewClassDesc
	nullReference byte
	PrevObject    *PrevObject
}

func ParseClassDesc(parser *StructuresParser) (*ClassDesc, error) {
	peekByte, err := parser.byteReader.PeekByte()
	if err != nil {
		return nil, err
	}
	classDesc := new(ClassDesc)
	switch peekByte {
	case TC_CLASSDESC:
		classDesc.newClassDesc, err = ParseNewClassDesc(parser)
		if err != nil {
			return nil, fmt.Errorf("In ParseClassDesc:\n %w", err)
		}
	case TC_NULL:
		classDesc.nullReference = TC_NULL
	case TC_REFERENCE:

	default:
		return nil, fmt.Errorf("In ParseClassDesc:No matched byte")
	}
	return classDesc, nil
}
