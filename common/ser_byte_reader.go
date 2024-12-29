package common

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type SerByteReader struct {
	Data        []byte
	Index       *int
	BytesWriten []byte
}

func (r *SerByteReader) ReadByte() (byte, error) {
	if *r.Index >= len(r.Data) {
		return 0, fmt.Errorf("end of Data reached")
	}
	theByte := r.Data[*r.Index]
	*r.Index++
	return theByte, nil

}

func (r *SerByteReader) ReadNByte(n int) ([]byte, error) {

	if *r.Index+n > len(r.Data) {
		return nil, fmt.Errorf("not enough Data to read")
	}
	result := make([]byte, n)
	copy(result, r.Data[*r.Index:*r.Index+n])
	*r.Index += n
	return result, nil
}

func (r *SerByteReader) PeekByte() (byte, error) {
	if *r.Index >= len(r.Data) {
		return 0, fmt.Errorf("end of Data reached")
	}
	theByte := r.Data[*r.Index]
	return theByte, nil
}

func (r *SerByteReader) PeekNByte(n int) ([]byte, error) {
	if *r.Index >= len(r.Data) {
		return nil, fmt.Errorf("end of Data reached")
	}
	theBytes := make([]byte, n)
	for i := 0; i < n; i++ {
		theBytes[i] = r.Data[*r.Index+i]
	}
	return theBytes, nil
}

func (r *SerByteReader) JumpByte() error {
	if *r.Index >= len(r.Data) {
		return fmt.Errorf("end of Data reached")
	}
	*r.Index++
	return nil
}

func (r *SerByteReader) ReadInt16() (int16, error) {
	byteArray, err := r.ReadNByte(2)
	if err != nil {
		return 0, err
	}
	result := int16(binary.BigEndian.Uint16(byteArray))
	return result, nil
}

// INT32 = JAVA INT
func (r *SerByteReader) ReadInt32() (int32, error) {
	byteArray, err := r.ReadNByte(4)
	if err != nil {
		return 0, err
	}
	result := int32(binary.BigEndian.Uint32(byteArray))
	return result, nil
}

func (r *SerByteReader) ReadLong() (int64, error) {
	byteArray, err := r.ReadNByte(8)
	if err != nil {
		return 0, nil
	}
	number := int64(binary.BigEndian.Uint64(byteArray))
	return number, nil
}

func (r *SerByteReader) ReadAllBytes() bool {
	return *r.Index == len(r.Data)
}

func (r *SerByteReader) WriteByte(b byte) error {
	r.BytesWriten = append(r.BytesWriten, b)
	return nil
}

func (r *SerByteReader) WriteNByte(bs []byte) {
	r.BytesWriten = append(r.BytesWriten, bs...)
}

func (r *SerByteReader) WriteNumber(num interface{}) error {
	buf := new(bytes.Buffer)
	switch num.(type) {
	case int8, int16, int32, int64,
		uint8, uint16, uint32, uint64,
		float32, float64:
		// 使用大端字节序写入缓冲区(与JAVA相同)
		err := binary.Write(buf, binary.BigEndian, num)
		if err != nil {
			return fmt.Errorf("failed to write value to buffer: %w", err)
		}
	default:
		return errors.New("unsupported type: only int/float types are allowed")
	}
	r.WriteNByte(buf.Bytes())
	return nil
}

func ConvertNumberToBytes(num interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	switch num.(type) {
	case int8, int16, int32, int64,
		uint8, uint16, uint32, uint64,
		float32, float64:
		// 使用大端字节序写入缓冲区(与JAVA相同)
		err := binary.Write(buf, binary.BigEndian, num)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("no supported type")
	}

	return buf.Bytes(), nil
}
