package common

import "fmt"

type SerUtfReader struct {
	NumBytesForUTF8 int
}

func (s *SerUtfReader) ReadUtf(datas []byte) (string, error) {
	if s.NumBytesForUTF8 == 1 {
		return string(datas), nil
	} else if s.NumBytesForUTF8 == 2 {
		return ReadAndUpdateWithTwoBytes()
	} else if s.NumBytesForUTF8 == 3 {
		return ReadAndUpdateWithThreeBytes()
	} else {
		return "", fmt.Errorf("wtf?")
	}
}

func ReadAndUpdateWithTwoBytes() (string, error)

func ReadAndUpdateWithThreeBytes() (string, error)
