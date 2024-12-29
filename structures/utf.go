package structures

import (
	"fmt"
	"main/common"
)

type UTF struct {
	Length    int16
	Value     string
	ByteArray []byte
}

type LongUTF struct {
	Length    int64
	Value     string
	ByteArray []byte
}

func NewUTF(length int16, value string, byteArray []byte) *UTF {
	return &UTF{
		Length:    length,
		Value:     value,
		ByteArray: byteArray,
	}
}

func ParseUtf(parser *StructuresParser) (*UTF, error) {
	length, value, byteArray, err := parser.ReadUtf()
	if err != nil {
		return nil, fmt.Errorf("in ParseUtf:\n %v", err)
	}
	return NewUTF(length, value, byteArray), nil
}

func (fn *UTF) ToByte(parser *StructuresParser) error {
	parser.utfReader.WriteUtf(fn.Value, fn.Length, &parser.ByteReader)
	return nil
}

func (u *UTF) ToString(indent int) (string, error) {
	var (
		result string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @Utf")
	indent += IndentSpaceCount
	byteArray, err := common.ConvertNumberToBytes(u.Length)
	if err != nil {
		return "", fmt.Errorf("in Utf#ToString:\n%v", err)
	}
	result += sb.Buildf("- length  ", []interface{}{u.Length, "  -  ", byteArray})
	result += sb.Buildf("- value  ", []interface{}{u.Value, "  -  ", u.ByteArray})
	return result, nil
}

func NewLongUTF(length int64, value string, byteArray []byte) *LongUTF {
	return &LongUTF{
		Length:    length,
		Value:     value,
		ByteArray: byteArray,
	}
}

func ParseLongUTF(parser *StructuresParser) (*LongUTF, error) {
	length, value, byteArray, err := parser.ReadLongUtf()
	if err != nil {
		return nil, fmt.Errorf("in ParserLongUTF: \n %v", err)
	}
	return NewLongUTF(length, value, byteArray), nil
}

func (ln *LongUTF) ToByte(parser *StructuresParser) error {
	parser.utfReader.WriteUtf(ln.Value, ln.Length, &parser.ByteReader)
	return nil
}

func (ln *LongUTF) ToString(indent int) (string, error) {
	var (
		result string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @LongUtf")
	indent += IndentSpaceCount
	byteArray, err := common.ConvertNumberToBytes(ln.Length)
	if err != nil {
		return "", fmt.Errorf("in LongUTF#ToString:\n%v", err)
	}
	result += sb.Buildf("- length  ", []interface{}{ln.Length, "  -  ", byteArray})
	result += sb.Buildf("- value  ", []interface{}{ln.Value, "  -  ", ln.ByteArray})
	return result, nil
}
