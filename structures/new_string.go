package structures

import (
	"fmt"
	"main/common"
)

type NewString struct {
	FLAG          byte
	TC_STRING     byte
	TC_LONGSTRING byte
	NewHandle     uint32
	UTF           *UTF
	LongUTF       *LongUTF
	Value         string
	PrevObject    *PrevObject
}

func NewNewString(handle uint32, utf *UTF) *NewString {
	return &NewString{
		FLAG:      TC_STRING,
		TC_STRING: TC_STRING,
		Value:     utf.Value,
		NewHandle: handle,
		UTF:       utf,
	}
}

func NewRefNewString(prevObject *PrevObject) *NewString {
	return &NewString{
		FLAG:       TC_REFERENCE,
		PrevObject: prevObject,
	}
}

func NewNewLongString(handle uint32, longUTF *LongUTF) *NewString {
	return &NewString{
		FLAG:          TC_LONGSTRING,
		TC_LONGSTRING: TC_LONGSTRING,
		Value:         longUTF.Value,
		NewHandle:     handle,
		LongUTF:       longUTF,
	}
}

func ParseNewString(parser *StructuresParser) (*NewString, error) {
	peekByte, err := parser.ByteReader.PeekByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseNewString:\n %v", err)
	}
	if peekByte == TC_REFERENCE {
		prevObject, err := ParsePrevObject(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseNewString:\n %v", err)
		}
		newString := NewRefNewString(prevObject)
		return newString, nil
	} else if peekByte == TC_STRING {
		parser.ByteReader.JumpByte()
		utf, err := ParseUtf(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseNewString:\n %v", err)
		}
		handle := parser.AddHandle()
		newString := NewNewString(handle, utf)
		parser.AddReferenced(newString)
		return newString, nil
	} else if peekByte == TC_LONGSTRING {
		parser.ByteReader.JumpByte()
		longUtf, err := ParseLongUTF(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseNewString:\n %v", err)
		}
		handle := parser.AddHandle()
		newString := NewNewLongString(handle, longUtf)
		parser.AddReferenced(newString)
		return newString, nil
	}
	return nil, fmt.Errorf("in ParseNewString: No TC_STRING OR TC_REFERENCE")
}

func (ns *NewString) ToByte(parser *StructuresParser) error {
	if ns.TC_STRING != 0 {
		parser.ByteReader.WriteByte(TC_STRING)
		ns.UTF.ToByte(parser)
		return nil
	} else if ns.TC_LONGSTRING != 0 {
		parser.ByteReader.WriteByte(TC_LONGSTRING)
		ns.LongUTF.ToByte(parser)
		return nil
	} else if ns.PrevObject != nil {
		err := ns.PrevObject.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in ParseNewString:\n %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("in WriteNewString:No TC_STRING,TC_LONGSTRING,PrevObject in the NewString")
	}
}

func (ns *NewString) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @NewString")
	indent += IndentSpaceCount
	switch ns.FLAG {
	case TC_STRING:
		result += sb.Buildf("- TC_STRING  ", []interface{}{TC_STRING})
		result += sb.Buildf("- newHandle  ", []interface{}{ns.NewHandle})
		str, err := ns.UTF.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in NewString#ToString:\n%v", err)
		}
		result += str

	case TC_LONGSTRING:
		result += sb.Buildf("- TC_LONGSTRING  ", []interface{}{TC_LONGSTRING})
		result += sb.Buildf("- newHandle  ", []interface{}{ns.NewHandle})
		str, err := ns.LongUTF.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in NewString#ToString:\n%v", err)
		}
		result += str
	case TC_REFERENCE:
		result += sb.Buildf("- TC_REFERENCE  ", []interface{}{TC_REFERENCE})
		str, err = ns.PrevObject.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in NewString#ToString:\n%v", err)
		}
		result += str
	}
	return result, nil
}

func (ns *NewString) GetNewHandle() uint32 {
	return ns.NewHandle
}
