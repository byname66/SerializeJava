package structures

import (
	"bytes"
	"fmt"
	"io"
	"main/common"
)

type Stream struct {
	STREAM_MAGIC   []byte
	STREAM_VERSION []byte
	Contents       []*Content
	SerVersionUIDs []SerVersionUID
	Num_TC_RESET   int
}

func NewStream(Contents []*Content, SerVersionUIDs []SerVersionUID, num int) *Stream {
	return &Stream{
		STREAM_MAGIC:   STREAM_MAGIC,
		STREAM_VERSION: STREAM_VERSION,
		Contents:       Contents,
		SerVersionUIDs: SerVersionUIDs,
		Num_TC_RESET:   num,
	}
}

func NewEmptyStream() *Stream {
	return &Stream{
		STREAM_MAGIC:   STREAM_MAGIC,
		STREAM_VERSION: STREAM_VERSION,
		Contents:       nil,
	}
}

// first to check the Magic and Version bytes
func CheckMagicAndVersion(parser *StructuresParser) error {
	tbytes, error := parser.ByteReader.ReadNByte(2)
	if error != nil {
		return fmt.Errorf("error reading Magic bytes: %v", error)
	}
	if !bytes.Equal(tbytes, STREAM_MAGIC) {
		return fmt.Errorf("magic bytes do not match")
	}
	tbytes, error = parser.ByteReader.ReadNByte(2)
	if error != nil {
		return fmt.Errorf("error reading Version bytes: %v", error)
	}
	if !bytes.Equal(tbytes, STREAM_VERSION) {
		return fmt.Errorf("version bytes do not match")
	}
	return nil
}

func ParseStream(parser *StructuresParser) (*Stream, error) {
	err := CheckMagicAndVersion(parser)
	if err != nil {
		return nil, err
	}
	contents := new([]*Content)
	num_TC_Reset := 0
	for {
		signByte, _ := parser.ByteReader.PeekByte()
		if signByte == TC_RESET {
			parser.ByteReader.ReadByte()
			num_TC_Reset += 1
		}
		if signByte != TC_RESET {
			break
		}
	}
	for {
		content, err := ParseContent(parser)
		if err != nil {
			if err == io.EOF {
				return NewStream(*contents, parser.SerVersionUIDs, num_TC_Reset), nil
			}
			return nil, err
		}

		*contents = append(*contents, content)
	}
}

func (stream *Stream) ToByte(numBytesForUtf8 int) (*StructuresParser, error) {
	parser := NewStructureParser(nil, numBytesForUtf8)
	parser.ByteReader.WriteNByte(STREAM_MAGIC)
	parser.ByteReader.WriteNByte(STREAM_VERSION)
	for i := 0; i < len(stream.Contents); i++ {
		stream.Contents[i].ToByte(parser)
	}
	return parser, nil
}

func (stream *Stream) ToString(indent int) (string, error) {
	sb := common.NewStringBuilder(&indent)
	result := sb.Buildf("- MAGIC:  ", []interface{}{stream.STREAM_MAGIC})
	result += sb.Buildf("- VERSION:  ", []interface{}{stream.STREAM_VERSION})
	result += sb.Build(" @Contents:")
	if stream.Num_TC_RESET != 0 {
		result += sb.BuildWithSpaces(fmt.Sprintf("@Content[0] - @Content[%v]  <TC_RESET>", stream.Num_TC_RESET), IndentSpaceCount)
		result += sb.BuildWithSpaces(fmt.Sprintf("- TC_RESET  %v", []interface{}{TC_RESET}), IndentSpaceCount*2)
	}
	if len(stream.Contents) == 1 {
		if stream.Num_TC_RESET != 0 {
			result += sb.BuildWithSpaces(fmt.Sprintf("@Content[%v]", stream.Num_TC_RESET+1), IndentSpaceCount)
		}
		str, err := stream.Contents[0].ToString(IndentSpaceCount)
		if err != nil {
			return "", err
		}
		result += str
	} else {
		for i := 0; i < len(stream.Contents); i++ {
			result += sb.BuildWithSpaces(fmt.Sprintf("@Content[%v]", i+stream.Num_TC_RESET), IndentSpaceCount)
			str, err := stream.Contents[i].ToString(IndentSpaceCount)
			if err != nil {
				return "", err
			}
			result += str
		}
	}
	return result, nil
}
