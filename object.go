package structures

import "fmt"

type Object struct {
	newObject     *NewObject
	newClass      *NewClass
	newArray      *NewArray
	newString     *NewString
	newEnum       *NewEnum
	newClassDesc  *NewClassDesc
	prevObject    *PrevObject
	nullReference *NullReference
	exception     *Exception
	tc_reset      byte
}

func ParseObject(parser *StructuresParser) (*Object, error) {
	peekByte, err := parser.byteReader.PeekByte()
	if err != nil {
		return nil, err
	}
	object := new(Object)
	switch peekByte {
	case TC_OBJECT:
		object.newObject, err = ParseNewObject(parser)
	case TC_CLASS:
		object.newClass, err = ParseNewClass(parser)
	case TC_ARRAY:
		object.newArray, err = ParseNewArray()
	case TC_STRING, TC_LONGSTRING:
		object.newString, err = ParseNewString(parser)
	case TC_ENUM:
		object.newEnum, err = ParseNewEnum()
	case TC_CLASSDESC:
		object.newClassDesc, err = ParseNewClassDesc(parser)
	case TC_REFERENCE:
		object.prevObject, err = ParsePreObject()
	case TC_NULL:
		object.nullReference, err = ParseNullReference()
	case TC_EXCEPTION:
		object.exception, err = ParseException()
	case TC_RESET:
		object.tc_reset = TC_RESET
	default:
		err = fmt.Errorf("In ParseObject: object not found")
	}
	if err != nil {
		return nil, err
	}
	return object, nil
}
