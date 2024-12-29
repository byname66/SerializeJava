package structures

import (
	"fmt"
	"main/common"
)

// This struct should be divided into 2 parts: ClassDesc and ProxyClassDesc.
type NewClassDesc struct {
	FLAG               byte
	TC_CLASSDESC       byte
	TC_PROXYCLASSDESC  byte
	ClassName          *UTF
	SerialVersionUID   int64
	NewHandle          uint32
	ClassDescInfo      *ClassDescInfo
	ProxyClassDescInfo *ProxyClassDescInfo
}

func NewNewClassDesc1(className *UTF, serialVersionUID int64, newHandle uint32, classDescInfo *ClassDescInfo, flag byte) *NewClassDesc {
	return &NewClassDesc{
		TC_CLASSDESC:     TC_CLASSDESC,
		ClassName:        className,
		SerialVersionUID: serialVersionUID,
		NewHandle:        newHandle,
		ClassDescInfo:    classDescInfo,
		FLAG:             TC_CLASSDESC,
	}
}

func NewNewClassDesc2(newHandle uint32, proxyClassDescInfo *ProxyClassDescInfo, flag byte) *NewClassDesc {
	return &NewClassDesc{
		TC_PROXYCLASSDESC:  TC_PROXYCLASSDESC,
		NewHandle:          newHandle,
		ProxyClassDescInfo: proxyClassDescInfo,
		FLAG:               TC_PROXYCLASSDESC,
	}
}

// Check do ClassDesc or ProxyClassDesc has super
func (ncd *NewClassDesc) HasSuper() bool {
	if ncd.ClassDescInfo != nil {
		return (ncd.ClassDescInfo.SuperClassDesc.TC_NULL != TC_NULL)
	} else {
		return (ncd.ProxyClassDescInfo.SuperClassDesc.TC_NULL != TC_NULL)
	}
}

func (ncd *NewClassDesc) GetSuper() *ClassDesc {
	if ncd.ClassDescInfo != nil {
		return ncd.ClassDescInfo.SuperClassDesc
	} else {
		return ncd.ProxyClassDescInfo.SuperClassDesc
	}
}

func ParseNewClassDesc(parser *StructuresParser) (*NewClassDesc, error) {
	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, err
	}
	var newClassDesc *NewClassDesc
	switch signByte {
	case TC_CLASSDESC:
		className, err := ParseUtf(parser)
		if err != nil {
			return nil, err
		}
		serialVersionUID, err := parser.ByteReader.ReadLong()
		if err != nil {
			return nil, err
		}
		newHandle := parser.AddHandle()
		classDescInfo, err := ParseClassDescInfo(parser)
		if err != nil {
			return nil, err
		}
		newClassDesc = NewNewClassDesc1(className, serialVersionUID, newHandle, classDescInfo, signByte)
		parser.StoreSerialVersionUID(className.Value, serialVersionUID, newClassDesc)
	case TC_PROXYCLASSDESC:
		newHandle := parser.AddHandle()
		proxyClassDescInfo, err := ParseProxyClassDescInfo(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseNewClassDesc:\n %v", err)
		}
		newClassDesc = NewNewClassDesc2(newHandle, proxyClassDescInfo, signByte)

	default:
		return nil, fmt.Errorf("newclassdesc not found")
	}

	parser.AddReferenced(newClassDesc)
	return newClassDesc, nil
}

func (newClassDesc *NewClassDesc) ToByte(parser *StructuresParser) error {
	if newClassDesc.TC_CLASSDESC != 0 {
		parser.ByteReader.WriteByte(TC_CLASSDESC)
		newClassDesc.ClassName.ToByte(parser)
		err := parser.ByteReader.WriteNumber(newClassDesc.SerialVersionUID)
		if err != nil {
			return fmt.Errorf("in WriteNewClassDesc:\n %v", err)
		}
		err = newClassDesc.ClassDescInfo.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteNewClassDesc:\n %v", err)
		}
		return nil
	} else if newClassDesc.TC_PROXYCLASSDESC != 0 {
		parser.ByteReader.WriteByte(TC_PROXYCLASSDESC)
		err := newClassDesc.ProxyClassDescInfo.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteNewClassDesc:\n %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("in WriteNewClassDesc:All content field are not exist")
	}
}

func (ncd *NewClassDesc) ToString(indent int) (string, error) {
	sb := common.NewStringBuilder(&indent)
	var (
		result string
		str    string
		err    error
	)
	if ncd.FLAG == TC_CLASSDESC {
		result += sb.Build(" @NewClassDesc")
		indent += IndentSpaceCount
		result += sb.Buildf("- TC_CLASSDESC  ", []interface{}{TC_CLASSDESC})
		result += sb.Build("@ClassName")
		res, err := ncd.ClassName.ToString(indent + IndentSpaceCount)
		if err != nil {
			return "", fmt.Errorf("in NewClassDesc#ToString:%v", err)
		}
		result += res
		byteArray, err := common.ConvertNumberToBytes(ncd.SerialVersionUID)
		if err != nil {
			return "", fmt.Errorf("in NewClassDesc#ToString:%v", err)
		}
		result += sb.Buildf("- SerialVersionUID  ", []interface{}{ncd.SerialVersionUID, "  -  ", byteArray})
		result += sb.Buildf("- newHandle  ", []interface{}{ncd.NewHandle})
		str, err = ncd.ClassDescInfo.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in NewClassDesc#ToString:%v", err)
		}
		result += str
	} else if ncd.FLAG == TC_PROXYCLASSDESC {
		result += sb.Build(" @ProxyClassDesc")
		indent += IndentSpaceCount
		result += sb.Buildf("- TC_PROXYCLASSDESC ", []interface{}{TC_PROXYCLASSDESC})
		result += sb.Buildf("- NewHandle  ", []interface{}{ncd.NewHandle})
		str, err = ncd.ProxyClassDescInfo.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in NewClassDesc#ToString:%v", err)
		}
		result += str
	}

	return result, nil
}

func (ncd *NewClassDesc) GetNewHandle() uint32 {
	return ncd.NewHandle
}
