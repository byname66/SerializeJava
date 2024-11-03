package common

import (
	"encoding/binary"
	"fmt"
)

type SerByteReader struct {
	Data  []byte
	Index *int
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

	if *r.Index+n >= len(r.Data) {
		return nil, fmt.Errorf("not enough Data to read") // 错误处理：数据不足
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

func (r *SerByteReader) ReadLong() (int64, error) {
	byteArray, err := r.ReadNByte(8)
	if err != nil {
		return 0, nil
	}
	number := int64(binary.BigEndian.Uint64(byteArray))
	return number, nil
}
