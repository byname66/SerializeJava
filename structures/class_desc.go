package structures

import (
	"fmt"
	"main/common"
)

type ClassDesc struct {
	NewClassDesc *NewClassDesc
	TC_NULL      byte
	PrevObject   *PrevObject
	Flag         byte
}

func ParseClassDesc(parser *StructuresParser) (*ClassDesc, error) {
	peekByte, err := parser.ByteReader.PeekByte()
	if err != nil {
		return nil, err
	}
	classDesc := new(ClassDesc)
	classDesc.Flag = peekByte
	switch peekByte {
	case TC_CLASSDESC, TC_PROXYCLASSDESC:
		classDesc.NewClassDesc, err = ParseNewClassDesc(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseClassDesc:\n %w", err)
		}
	case TC_NULL:
		classDesc.TC_NULL = TC_NULL
		parser.ByteReader.JumpByte()
	case TC_REFERENCE:
		prevObject, err := ParsePrevObject(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseClassDesc:\n %w", err)
		}
		classDesc.PrevObject = prevObject
	default:
		return nil, fmt.Errorf("in ParseClassDesc:No matched byte")
	}
	return classDesc, nil
}

func (classDesc *ClassDesc) ToByte(parser *StructuresParser) error {
	if classDesc.NewClassDesc != nil {
		err := classDesc.NewClassDesc.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteClassDesc:\n %v", err)
		}
		return nil
	} else if classDesc.TC_NULL != 0 {
		parser.ByteReader.WriteByte(TC_NULL)
		return nil
	} else if classDesc.PrevObject != nil {
		err := classDesc.PrevObject.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteClassDesc:\n %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("in WriteClassDesc:All content field are not exist")
	}
}

func (classDesc *ClassDesc) ToString(indent int) (string, error) {
	sb := common.NewStringBuilder(&indent)
	var (
		result string
		str    string
		err    error
	)
	result += sb.Build(" @ClassDesc")
	indent += IndentSpaceCount
	flag := classDesc.Flag
	switch flag {
	case TC_CLASSDESC, TC_PROXYCLASSDESC:
		str, err = classDesc.NewClassDesc.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in ClassDesc#ToString:\n%v", err)
		}
		result += str
	case TC_NULL:
		result += sb.Buildf("- TC_NULL  ", []interface{}{classDesc.TC_NULL})
	case TC_REFERENCE:
		str, err = classDesc.PrevObject.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in ClassDesc#ToString:\n%v", err)
		}
		result += str
	}
	return result, nil
}
