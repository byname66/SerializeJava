package structures

import (
	"encoding/binary"
	"fmt"
	"main/common"
	"strings"
)

type NewArray struct {
	TypeCode  string
	TC_ARRAY  byte
	ClassDesc *ClassDesc
	NewHandle uint32
	Size      int32
	Values    []*Value
}

func NewNewArray(classDesc *ClassDesc, newHandle uint32, size int32, values []*Value, typeCode string) *NewArray {
	return &NewArray{
		TC_ARRAY:  TC_ARRAY,
		ClassDesc: classDesc,
		NewHandle: newHandle,
		Size:      size,
		Values:    values,
		TypeCode:  typeCode,
	}
}

func ParseNewArray(parser *StructuresParser) (*NewArray, error) {
	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseNewArray:\n %v", err)
	}
	if signByte != TC_ARRAY {
		return nil, fmt.Errorf("in ParseNewArray: No TC_ARRAY")
	}
	classDesc, err := ParseClassDesc(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseNewArray:\n %v", err)
	}
	newHandle := parser.AddHandle()
	byteArray, err := parser.ByteReader.ReadNByte(4)
	if err != nil {
		return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
	}
	size := int32(binary.BigEndian.Uint32(byteArray))
	var values []*Value
	var name string
	if classDesc.Flag == TC_CLASSDESC {
		name = classDesc.NewClassDesc.ClassName.Value
	} else if classDesc.Flag == TC_REFERENCE {
		obj, err := parser.GetReferenced(classDesc.PrevObject.Handler)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		if des, ok := obj.(*NewClassDesc); ok {
			if des.ClassDescInfo != nil {
				name = des.ClassName.Value
			} else if des.ProxyClassDescInfo != nil {
				return nil, fmt.Errorf("not allow ProxyClassDesc")
			}
		} else {
			return nil, fmt.Errorf("only NewClassDesc here")
		}
	} else if classDesc.Flag == TC_PROXYCLASSDESC || classDesc.Flag == TC_NULL {
		return nil, fmt.Errorf("not allow TC_NULL or TC_PROXYCLASSDESC")
	}
	if !strings.HasPrefix(name, "[") || len(name) < 2 {
		return nil, fmt.Errorf("JAVA_TC_ARRAY ClassName %v", name)
	}
	type_code := name[1:2]
	for i := 0; i < int(size); i++ {
		value, err := ParseValue(parser, type_code)
		if err != nil {
			return nil, fmt.Errorf("in ParseNewArray:\n %v", err)
		}
		values = append(values, value)
	}
	newArray := NewNewArray(classDesc, newHandle, size, values, type_code)
	parser.AddReferenced(newArray)
	return newArray, nil
}

func (na *NewArray) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteByte(TC_ARRAY)
	err := na.ClassDesc.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteNewArray:\n %v", err)
	}

	parser.ByteReader.WriteNumber(na.Size)

	for i := 0; i < len(na.Values); i++ {
		err = na.Values[i].ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteNewArray:\n %v", err)
		}
	}

	return nil
}

func (na *NewArray) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @NewArray")
	result += sb.Buildf("- TC_ARRAY  ", []interface{}{TC_ARRAY})
	str, err = na.ClassDesc.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in NewArray#ToString\n%v", err)
	}
	result += str
	result += sb.Buildf("- newHandle  ", []interface{}{na.NewHandle})
	byteArray, err := common.ConvertNumberToBytes(na.Size)
	if err != nil {
		return "", fmt.Errorf("in NewArray#ToString\n%v", err)
	}
	result += sb.Buildf("- size  ", []interface{}{na.Size, "  -  ", byteArray})
	result += sb.Build(" @Values")
	indent += IndexToArraySpaceCount
	for i := 0; i < len(na.Values); i++ {
		result += sb.Buildf("Index  [", []interface{}{i, "]"})
		str, err = na.Values[i].ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in NewArray#ToString\n%v", err)
		}
		result += str
	}
	return result, nil
}

func (na *NewArray) GetNewHandle() uint32 {
	return na.NewHandle
}
