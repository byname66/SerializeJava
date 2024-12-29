package common

import (
	"fmt"
)

type SerUtfReader struct {
	NumBytesForUTF8 int
}

func (s *SerUtfReader) ReadUtf(n int, serByteReader *SerByteReader) (string, []byte, error) {
	num := n * s.NumBytesForUTF8
	byteArray, err := serByteReader.ReadNByte(num)
	if err != nil {
		return "", nil, fmt.Errorf("in ReadNUtf:\n %v", err)
	}
	str := string(byteArray)
	return str, byteArray, nil
}

func (s *SerUtfReader) WriteUtf(str string, length interface{}, byteReader *SerByteReader) error {
	if s.NumBytesForUTF8 == 1 {
		result := []byte(str)
		byteReader.WriteNumber(length)
		byteReader.WriteNByte(result)
	} else if s.NumBytesForUTF8 == 2 {
		err := write1With2Byte(str, length, byteReader)
		if err != nil {
			return fmt.Errorf("in Write1With2Byte: %s", err)
		}
	} else if s.NumBytesForUTF8 == 3 {
		err := write1With3Byte(str, length, byteReader)
		return fmt.Errorf("in Write1With3Byte: %s", err)
	}
	return nil
}

func write1With2Byte(str string, length interface{}, byteReader *SerByteReader) error {
	var bytes []byte
	var dataLen int
	switch v := length.(type) {
	case int16:
		dataLen = int(v) * 2
		byteReader.WriteNumber(int16(dataLen))
	case int64:
		dataLen = int(v) * 2
		byteReader.WriteNumber(int64(dataLen))
	default:
		return fmt.Errorf("unsupported length type: %T", length)
	}
	bytes = make([]byte, dataLen)
	for i, j := 0, 0; i < dataLen/2 && j < dataLen; i, j = i+1, j+2 {
		if i >= len(str) {
			return fmt.Errorf("string index out of range: %d", i)
		}
		ch := string(str[i])
		datas, ok := TowFor1Table[ch]
		if !ok || len(datas) < 2 {
			return fmt.Errorf("invalid mapping for character: %s", ch)
		}
		bytes[j] = byte(datas[0])
		bytes[j+1] = byte(datas[1])
	}
	byteReader.WriteNByte(bytes)
	return nil
}

func write1With3Byte(str string, length interface{}, byteReader *SerByteReader) error {
	var bytes []byte
	var dataLen int
	switch v := length.(type) {
	case int16:
		dataLen = int(v) * 3
		byteReader.WriteNumber(int16(dataLen))
	case int64:
		dataLen = int(v) * 3
		byteReader.WriteNumber(int64(dataLen))
	default:
		return fmt.Errorf("unsupported length type: %T", length)
	}
	bytes = make([]byte, dataLen)
	for i, j := 0, 0; i < dataLen/3 && j < dataLen; i, j = i+1, j+3 {
		if i >= len(str) {
			return fmt.Errorf("string index out of range: %d", i)
		}
		ch := string(str[i])
		datas, ok := ThreeFor1Table[ch]
		if !ok || len(datas) < 2 {
			return fmt.Errorf("invalid mapping for character: %s", ch)
		}
		bytes[j] = byte(datas[0])
		bytes[j+1] = byte(datas[1])
		bytes[j+2] = byte(datas[2])
	}
	byteReader.WriteNByte(bytes)
	return nil
}
