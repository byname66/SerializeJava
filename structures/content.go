package structures

import (
	"fmt"
	"io"
)

type Content struct {
	BlockData *BlockData
	Object    *Object
}

func ParseContent(parser *StructuresParser) (*Content, error) {
	if parser.ByteReader.ReadAllBytes() {
		return nil, io.EOF
	}
	peekByte, err := parser.ByteReader.PeekByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseContent:\n %v", err)
	}
	content := new(Content)
	switch peekByte {
	case TC_BLOCKDATA, TC_BLOCKDATALONG:
		content.BlockData, err = ParseBlockData(parser)
	case TC_BLOCKDATA, TC_BLOCKDATALONG, TC_OBJECT, TC_CLASS, TC_ARRAY, TC_STRING, TC_LONGSTRING, TC_ENUM, TC_CLASSDESC, TC_REFERENCE, TC_NULL, TC_EXCEPTION, TC_RESET:
		content.Object, err = ParseObject(parser)
	default:
		return nil, fmt.Errorf("in ParseContent: No Object or BlockData byte")
	}
	if err != nil {
		return nil, fmt.Errorf("in ParseContent:\n %v", err)
	}
	return content, nil
}

func (content *Content) ToByte(parser *StructuresParser) error {
	if content.Object != nil {
		content.Object.ToByte(parser)
		return nil
	} else if content.BlockData != nil {
		content.BlockData.ToByte(parser)
		return nil
	} else {
		return fmt.Errorf("in WriteContent:All content field are nil")
	}
}

func (content *Content) ToString(indent int) (string, error) {
	var result string
	var err error
	if content.Object != nil {
		result, err = content.Object.ToString(indent + IndentSpaceCount)
		if err != nil {
			return "", fmt.Errorf("in content#ToString: \n%v", err)
		}
	} else if content.BlockData != nil {
		result, err = content.BlockData.ToString(indent + IndentSpaceCount)
		if err != nil {
			return "", fmt.Errorf("in content#ToString: \n%v", err)
		}
	} else {
		return "", fmt.Errorf("in content#ToString: No Object or BlockData")
	}
	return result, nil
}
