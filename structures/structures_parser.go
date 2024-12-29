package structures

import (
	"encoding/base64"
	"fmt"
	"io"
	"main/common"
	"os"
	"strings"
)

var index int

type ReferencedObject interface {
	GetNewHandle() uint32
}

type SerVersionUID struct {
	ClassName        string
	SerialVersionUID int64
	StructPtr        *NewClassDesc
}

type StructuresParser struct {
	ByteReader        common.SerByteReader
	utfReader         common.SerUtfReader
	index             int
	newHandle         uint32
	referencedObjects map[uint32]ReferencedObject
	SerVersionUIDs    []SerVersionUID
}

func NewStructureParser(Data []byte, numBytesForUtf8 int) *StructuresParser {
	index = 0
	return &StructuresParser{
		index:     index,
		newHandle: BASE_WRITE_HANDLE,
		ByteReader: common.SerByteReader{
			Data:  Data,
			Index: &index,
		},
		utfReader: common.SerUtfReader{
			NumBytesForUTF8: numBytesForUtf8,
		},
		referencedObjects: make(map[uint32]ReferencedObject),
		SerVersionUIDs:    make([]SerVersionUID, 0),
	}
}

func Parse(Data []byte, numBytesForUtf8 int) (*Stream, error) {

	parser := NewStructureParser(Data, numBytesForUtf8)
	stream, err := ParseStream(parser)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func ParseBase64Data(base64Data string, numBytesForUtf8 int) (*Stream, error) {
	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}
	return Parse(decodedData, numBytesForUtf8)
}

func ParseBase64FileData(filename string, numBytesForUtf8 int) (*Stream, error) {
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
	return p.newHandle - 1
}

func (p *StructuresParser) AddReferenced(obj ReferencedObject) {
	p.referencedObjects[obj.GetNewHandle()] = obj
}

func (p *StructuresParser) GetReferenced(handle uint32) (ReferencedObject, error) {
	return p.referencedObjects[handle], nil
}

// Pass a *ClassDesc as input, encapsulating its SuperClassDesc(*ClassDesc) and the *ClassDesc referenced by it into the result and return.
func (p *StructuresParser) FindAllDescs(classDesc *ClassDesc) ([]*ClassDesc, error) {
	var result []*ClassDesc
	var superResult []*ClassDesc
	desc := new(ClassDesc)
	super := new(ClassDesc)
	var err error
	//First, get the referenced NewClassDesc and assign it to a ClassDesc
	if classDesc.Flag == TC_CLASSDESC || classDesc.Flag == TC_PROXYCLASSDESC {
		desc = classDesc
		if classDesc.Flag == TC_CLASSDESC {
			result = append(result, desc)
		}
	} else if classDesc.Flag == TC_REFERENCE {
		obj, err := p.GetReferenced(classDesc.PrevObject.Handler)
		if err != nil {
			return nil, err
		}
		switch ncd := obj.(type) {
		case *NewClassDesc:
			desc.NewClassDesc = ncd
		default:
			return nil, fmt.Errorf("unsupportType:Only NewClassDesc(ProxyCLassDesc or ClassDesc)")
		}
		if desc.NewClassDesc.ClassDescInfo != nil {
			result = append(result, desc)
		}
	} else {
		return nil, fmt.Errorf("null ClassDesc")
	}
	//Then, check if it has a super ClassDesc. If it does, recursively call FindAllDescs to find all the super ClassDesc objects.
	if desc.NewClassDesc.HasSuper() {
		super = desc.NewClassDesc.GetSuper()
		superResult, err = p.FindAllDescs(super)
		if err != nil {
			return nil, err
		}
		result = append(result, superResult...)
	}

	return result, nil
}

// Store className and the class's serialVersionUID
func (p *StructuresParser) StoreSerialVersionUID(className string, UID int64, ptr *NewClassDesc) {
	obj := new(SerVersionUID)
	obj.ClassName = className
	obj.SerialVersionUID = UID
	obj.StructPtr = ptr
	p.SerVersionUIDs = append(p.SerVersionUIDs, *obj)
}

// If next block of bytes containes a name(fieldName or className),
// which always uses two byte to express its length,use this function
// to return the length(int16) and the name.
func (p *StructuresParser) ReadUtf() (int16, string, []byte, error) {
	bytes, err := p.ByteReader.PeekNByte(3)
	if err != nil {
		return 0, "", nil, err
	}
	tmp_length, err := p.ByteReader.ReadInt16()
	if err != nil {
		return 0, "", nil, err
	}
	if bytes[2] == 0x80 || bytes[2] == 0x81 { //If UTF OverLong Encoding(1 UTF for 2 byte)
		return p.Read1UtfWith2Bytes(tmp_length)
	} else if bytes[2] == 0xe0 {
		return p.Read1UtfWith3Bytes(tmp_length)
	} else {
		length := tmp_length
		res, byteArray, err := p.utfReader.ReadUtf(int(length), &p.ByteReader)
		if err != nil {
			return 0, "", nil, err
		}
		return length, res, byteArray, nil
	}
}

// Maybe something will wrong because of converting int64 to int
func (p *StructuresParser) ReadLongUtf() (int64, string, []byte, error) {
	length, err := p.ByteReader.ReadLong()
	if err != nil {
		return 0, "", nil, err
	}
	res, byteArray, err := p.utfReader.ReadUtf(int(length), &p.ByteReader)
	if err != nil {
		return 0, "", nil, err
	}
	return length, res, byteArray, nil
}

func (p *StructuresParser) Read1UtfWith2Bytes(len int16) (int16, string, []byte, error) {
	var builder strings.Builder
	bytes, err := p.ByteReader.ReadNByte(int(len))
	if err != nil {
		return 0, "", nil, fmt.Errorf("in Read1UtfWith2Bytes:\n %v", err)
	}
	res_len := len / 2
	builder.Grow(int(res_len))
	for i := 0; i < int(len); i += 2 {
		b1 := int(bytes[i])
		b2 := int(bytes[i+1])
		builder.WriteString(string(((b1 & 0x1F) << 6) | ((b2 & 0x3F) << 0)))
	}
	return res_len, builder.String(), bytes, nil
}

func (p *StructuresParser) Read1UtfWith3Bytes(len int16) (int16, string, []byte, error) {
	var builder strings.Builder
	bytes, err := p.ByteReader.ReadNByte(int(len))
	if err != nil {
		return 0, "", nil, fmt.Errorf("in Read1UtfWith2Bytes:\n %v", err)
	}
	res_len := len / 3
	builder.Grow(int(res_len))
	for i := 0; i < int(len); i += 3 {
		b1 := int(bytes[i])
		b2 := int(bytes[i+1])
		b3 := int(bytes[i+2])
		builder.WriteString(string((((b1 & 0x0F) << 12) | ((b2 & 0x3F) << 6) | ((b3 & 0x3F) << 0))))
	}
	return res_len, builder.String(), bytes, nil
}
