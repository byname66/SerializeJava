package structures

import (
	"fmt"
	"main/common"
)

type Object struct {
	FLAG         byte
	NewObject    *NewObject
	NewClass     *NewClass
	NewArray     *NewArray
	NewString    *NewString
	NewEnum      *NewEnum
	NewClassDesc *NewClassDesc
	PrevObject   *PrevObject
	TC_NULL      byte
	Exception    *Exception
	TC_RESET     byte
}

func ParseObject(parser *StructuresParser) (*Object, error) {
	peekByte, err := parser.ByteReader.PeekByte()
	if err != nil {
		return nil, err
	}
	object := new(Object)
	object.FLAG = peekByte
	switch peekByte {
	case TC_OBJECT:
		object.NewObject, err = ParseNewObject(parser)
	case TC_CLASS:
		object.NewClass, err = ParseNewClass(parser)
	case TC_ARRAY:
		object.NewArray, err = ParseNewArray(parser)
	case TC_STRING, TC_LONGSTRING:
		object.NewString, err = ParseNewString(parser)
	case TC_ENUM:
		object.NewEnum, err = ParseNewEnum(parser)
	case TC_CLASSDESC:
		object.NewClassDesc, err = ParseNewClassDesc(parser)
	case TC_REFERENCE:
		object.PrevObject, err = ParsePrevObject(parser)
	case TC_NULL:
		object.TC_NULL = TC_NULL
		parser.ByteReader.JumpByte()
	case TC_EXCEPTION:
		object.Exception, err = ParseException(parser)
	case TC_RESET:
		object.TC_RESET = TC_RESET
		parser.ByteReader.JumpByte()
	default:
		err = fmt.Errorf("in ParseObject: object not found")
	}
	if err != nil {
		return nil, err
	}
	return object, nil
}

// func (object *Object) ToByte(parser *StructuresParser) error {
// 	if object.NewObject != nil {
// 		err := object.NewObject.ToByte(parser)
// 		if err != nil {
// 			return fmt.Errorf("in WriteObject:\n %v", err)
// 		}
// 		return nil
// 	} else if object.NewClass != nil {
// 		err := object.NewClass.ToByte(parser)
// 		if err != nil {
// 			return fmt.Errorf("in WriteObject:\n %v", err)
// 		}
// 		return nil
// 	} else if object.NewArray != nil {
// 		err := object.NewArray.ToByte(parser)
// 		if err != nil {
// 			return fmt.Errorf("in WriteObject:\n %v", err)
// 		}
// 		return nil
// 	} else if object.NewString != nil {
// 		err := object.NewString.ToByte(parser)
// 		if err != nil {
// 			return fmt.Errorf("in WriteObject:\n %v", err)
// 		}
// 		return nil
// 	} else if object.NewEnum != nil {
// 		err := object.NewEnum.ToByte(parser)
// 		if err != nil {
// 			return fmt.Errorf("in WriteObject:\n %v", err)
// 		}
// 		return nil
// 	} else if object.NewClassDesc != nil {
// 		err := object.NewClassDesc.ToByte(parser)
// 		if err != nil {
// 			return fmt.Errorf("in WriteObject:\n %v", err)
// 		}
// 		return nil
// 	} else if object.PrevObject != nil {
// 		err := object.PrevObject.ToByte(parser)
// 		if err != nil {
// 			return fmt.Errorf("in WriteObject:\n %v", err)
// 		}
// 		return nil
// 	} else if object.TC_NULL != 0 {
// 		parser.ByteReader.WriteByte(object.TC_NULL)
// 		return nil
// 	} else if object.Exception != nil {
// 		err := object.Exception.ToByte(parser)
// 		if err != nil {
// 			return fmt.Errorf("in WriteObject:\n %v", err)
// 		}
// 		return nil
// 	} else if object.TC_RESET != 0 {
// 		parser.ByteReader.WriteByte(object.TC_RESET)
// 		return nil
// 	} else {
// 		return fmt.Errorf("in WriteObject:All content field are not exist")
// 	}

// }

func (object *Object) ToByte(parser *StructuresParser) error {
	Flag := object.FLAG
	var err error
	switch Flag {
	case TC_OBJECT:
		err = object.NewObject.ToByte(parser)
	case TC_CLASS:
		err = object.NewClass.ToByte(parser)
	case TC_ARRAY:
		err = object.NewArray.ToByte(parser)
	case TC_STRING, TC_LONGSTRING:
		err = object.NewString.ToByte(parser)
	case TC_ENUM:
		err = object.NewEnum.ToByte(parser)
	case TC_CLASSDESC:
		err = object.NewClassDesc.ToByte(parser)
	case TC_REFERENCE:
		err = object.PrevObject.ToByte(parser)
	case TC_NULL:
		parser.ByteReader.WriteByte(object.TC_NULL)
		err = nil
	case TC_EXCEPTION:
		err = object.Exception.ToByte(parser)
	case TC_RESET:
		parser.ByteReader.WriteByte(object.TC_RESET)
		err = nil
	default:
		err = fmt.Errorf("in ParseObject: object not found")
	}
	if err != nil {
		return err
	}
	return nil
}

func (object *Object) ToString(indent int) (string, error) {
	sb := common.NewStringBuilder(&indent)
	var result string
	var err error
	Flag := object.FLAG
	switch Flag {
	case TC_OBJECT:
		result, err = object.NewObject.ToString(indent)
	case TC_CLASS:
		result, err = object.NewClass.ToString(indent)
	case TC_ARRAY:
		result, err = object.NewArray.ToString(indent)
	case TC_STRING, TC_LONGSTRING:
		result, err = object.NewString.ToString(indent)
	case TC_ENUM:
		result, err = object.NewEnum.ToString(indent)
	case TC_CLASSDESC:
		result, err = object.NewClassDesc.ToString(indent)
	case TC_REFERENCE:
		result, err = object.PrevObject.ToString(indent)
	case TC_NULL:
		result = sb.Buildf("- TC_NULL  ", []interface{}{TC_NULL})
	case TC_EXCEPTION:
		result, err = object.Exception.ToString(indent)
	case TC_RESET:
		result = sb.Buildf("- TC_RESET  ", []interface{}{TC_RESET})
	default:
		result = ""
		err = fmt.Errorf("in object#ToString:FLAG wrong")
	}
	if err != nil {
		return "", fmt.Errorf("in object#ToString:\n%v", err)
	}
	return result, nil
}
