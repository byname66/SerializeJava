package structures

import "fmt"

type NewObject struct {
	tc_object byte
	classDesc *ClassDesc
	newHandle uint32
}

func ParseNewObject(parser *StructuresParser) (*NewObject, error) {
	signByte, err := parser.byteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("In ParseNewObject:\n %w", err)
	}
	if signByte != TC_OBJECT {
		return nil, fmt.Errorf("In ParseNewObject: No TC_OBJECT")
	}
	classDesc, err := ParseClassDesc(parser)
	if err != nil {
		return nil, fmt.Errorf("In ParseNewObject:\n %w", err)
	}
	newHandle := parser.AddHandle()

}
