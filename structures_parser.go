package structures

import (
	"encoding/base64"
	"io"
	"main/common"
	"os"
)

var index int

type StructuresParser struct {
	byteReader common.SerByteReader
	utfReader  common.SerUtfReader
	index      int
	newHandle  uint32
}

func NewStructureParser(Data []byte, numBytesForUtf8 int) *StructuresParser {
	index = 0
	return &StructuresParser{
		index:     index,
		newHandle: 8257535,
		byteReader: common.SerByteReader{
			Data:  Data,
			Index: &index,
		},
		utfReader: common.SerUtfReader{
			NumBytesForUTF8: numBytesForUtf8,
		},
	}
}

func Parse(Data []byte, numBytesForUtf8 int) (*Stream, error) {

	parser := NewStructureParser(Data, numBytesForUtf8)
	return ParseStream(parser)

}

func ParseBase64Data(base64Data string, numBytesForUtf8 int) (*Stream, error) {
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}
	return Parse(decodedData, numBytesForUtf8)
}

func ParseFileData(filename string, numBytesForUtf8 int) (*Stream, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	Data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	content := string(Data)
	return ParseBase64Data(content, numBytesForUtf8)

}

func (p *StructuresParser) AddHandle() uint32 {
	p.newHandle++
	return p.newHandle
}

// If next block of bytes containes a name(fieldName or className),
// which always uses two byte to express its length,use this function
// to return the length(int16) and the name.
func (p *StructuresParser) ReadName() (int16, string, error) {
	nameLength, err := p.byteReader.ReadInt16()
	if err != nil {
		return 0, "", err
	}
	nameBytes, err := p.byteReader.ReadNByte(int(nameLength))
	if err != nil {
		return 0, "", err
	}
	fieldName, err := p.utfReader.ReadUtf(nameBytes)
	if err != nil {
		return 0, "", err
	}
	return int16(nameLength), fieldName, nil
}
