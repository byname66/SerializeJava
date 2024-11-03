package structures

import "fmt"

type Content struct {
	blockData *BlockData
	object    *Object
}

func ParseContent(parser *StructuresParser) (*Content, error) {
	peekByte, err := parser.byteReader.PeekByte()
	if err != nil {
		return nil, fmt.Errorf("In ParseContent:\n %w", err)
	}
	content := new(Content)
	switch peekByte {
	case TC_BLOCKDATA, TC_BLOCKDATALONG:
		content.blockData, err = ParseBlockData()
	case TC_BLOCKDATA, TC_BLOCKDATALONG, TC_OBJECT, TC_CLASS, TC_ARRAY, TC_STRING, TC_LONGSTRING, TC_ENUM, TC_CLASSDESC, TC_REFERENCE, TC_NULL, TC_EXCEPTION, TC_RESET:
		content.object, err = ParseObject(parser)
	default:
		return nil, fmt.Errorf("In ParseContent: No Object or BlockData byte")
	}
	if err != nil {
		return nil, fmt.Errorf("In ParseContent:\n %w", err)
	}
	return content, nil
}
