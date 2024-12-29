package ui

import (
	"encoding/base64"
	"fmt"
	"main/structures"
	"strconv"
)

var stream *structures.Stream
var err error

func ConvertInputToStream(numBytesForUtf8 int) (*structures.Stream, error) {
	if IsBase64Valid(userinput) {
		stream, err = structures.ParseBase64Data(userinput, 1)
		if err != nil {
			return nil, err
		}
	} else {
		stream, err = structures.ParseBase64FileData(userinput, 1)
		if err != nil {
			return nil, err
		}
	}
	return stream, nil
}

func domain() {
	stream, err := ConvertInputToStream(1)
	if err != nil {
		result = err.Error()
		return
	}
	res, err := stream.ToString(0)
	if err != nil {
		result = err.Error()
	}
	result = res

}

// Insert lots of TC_RESET byte into the stream to bypass the WAF, which imposes a data length limitation.
func insertDirtyData(n string, stream *structures.Stream) (*structures.Stream, error) {
	intVal, err := strconv.Atoi(n)
	if err != nil {
		return nil, err
	}
	dirtyContent := new(structures.Content)
	dirtyObject := new(structures.Object)
	dirtyObject.TC_RESET = structures.TC_RESET
	dirtyObject.FLAG = structures.TC_RESET
	dirtyContent.Object = dirtyObject
	dirtyContents := make([]*structures.Content, intVal)
	for i := 0; i < intVal; i++ {
		dirtyContents[i] = dirtyContent
	}
	struContents := stream.Contents
	resContents := append(dirtyContents, struContents...)
	stream.Contents = resContents
	return stream, nil
}

func ConvertStreamToBase64(stream *structures.Stream, numBytesForUtf8 int) (string, error) {
	byteWriteParser, err := stream.ToByte(numBytesForUtf8)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	return base64.StdEncoding.EncodeToString(byteWriteParser.ByteReader.BytesWriten), nil
}

func IsBase64Valid(base64Str string) bool {
	_, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return false
	}
	return true
}
