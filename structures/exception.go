package structures

import (
	"fmt"
	"main/common"
)

type Exception struct {
	TC_EXCEPTION    byte
	TC_RESET_BEGIN  byte
	ResetObjects    []*Object
	TC_RESET_FINISH byte
}

func NewException(resetObjects []*Object) *Exception {
	return &Exception{
		TC_EXCEPTION:    TC_EXCEPTION,
		TC_RESET_BEGIN:  TC_RESET,
		ResetObjects:    resetObjects,
		TC_RESET_FINISH: TC_RESET,
	}
}

func ParseException(parser *StructuresParser) (*Exception, error) {
	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseException:\n %v", err)
	}
	if signByte != TC_EXCEPTION {
		return nil, fmt.Errorf("in ParseException: NO TC_EXCEPTION")
	}
	_, err = parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseException:\n %v", err)
	}
	var ResetObjects []*Object
	for {
		object, err := ParseObject(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseException:\n %v", err)
		}
		if object.TC_RESET == TC_RESET {
			break
		}
		ResetObjects = append(ResetObjects, object)
	}
	return NewException(ResetObjects), nil
}

func (ex *Exception) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteByte(TC_EXCEPTION)
	parser.ByteReader.WriteByte(ex.TC_RESET_BEGIN)
	for i := 0; i < len(ex.ResetObjects); i++ {
		err := ex.ResetObjects[i].ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteException:\n %v", err)
		}
	}
	parser.ByteReader.WriteByte(ex.TC_RESET_FINISH)
	return nil
}

func (ex *Exception) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @Exception")
	indent += IndentSpaceCount
	result += sb.Buildf("- TC_EXCEPTION", []interface{}{TC_EXCEPTION})
	result += sb.Buildf("- TC_RESET_BEGIN - TC_RESET", []interface{}{TC_RESET})
	result += sb.Build(" @ResetObjects")
	indent += IndexToArraySpaceCount
	for i := 0; i < len(ex.ResetObjects); i++ {
		result += sb.Buildf("Index  [", []interface{}{i, "]"})
		str, err = ex.ResetObjects[i].ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in Exception#ToString\n%v", err)
		}
		result += str
	}
	indent -= IndexToArraySpaceCount
	result += sb.Buildf("- TC_RESET_FINISH - TC_RESET", []interface{}{TC_RESET})
	return result, nil
}
