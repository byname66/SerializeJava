package structures

import (
	"encoding/binary"
	"fmt"
)

type NewClassDesc struct {
	tc_classdesc       byte
	tc_proxy           byte
	className          string
	serialVersionUID   int64
	newHandle          uint32
	classDescInfo      *ClassDescInfo
	proxyClassDescInfo *ProxyClassDescInfo
}

func NewNewClassDesc1(className string, serialVersionUID int64, newHandle uint32, classDescInfo *ClassDescInfo) *NewClassDesc {
	return &NewClassDesc{
		tc_classdesc:     TC_CLASSDESC,
		className:        className,
		serialVersionUID: serialVersionUID,
		newHandle:        newHandle,
		classDescInfo:    classDescInfo,
	}
}

func ParseNewClassDesc(parser *StructuresParser) (*NewClassDesc, error) {
	signByte, err := parser.byteReader.ReadByte()
	if err != nil {
		return nil, err
	}
	var newClassDesc *NewClassDesc
	switch signByte {
	case TC_CLASSDESC:
		byteArray, err := parser.byteReader.ReadNByte(2)
		classNameLength := int(binary.BigEndian.Uint16(byteArray))
		if err != nil {
			return nil, err
		}
		classNameBytes, err := parser.byteReader.ReadNByte(classNameLength)
		if err != nil {
			return nil, err
		}
		className, err := parser.utfReader.ReadUtf(classNameBytes)
		if err != nil {
			return nil, err
		}
		serialVersionUID, err := parser.byteReader.ReadLong()
		if err != nil {
			return nil, err
		}
		newHandle := parser.AddHandle()
		classDescInfo, err := ParseClassDescInfo(parser)
		if err != nil {
			return nil, err
		}
		newClassDesc = NewNewClassDesc1(className, serialVersionUID, newHandle, classDescInfo)

	case TC_PROXYCLASSDESC:

	default:
		err = fmt.Errorf("newclassdesc not found")
	}
	return newClassDesc, nil
}
