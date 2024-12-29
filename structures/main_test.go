package structures

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestEcho(t *testing.T) {
	file, error := os.Open("E:\\CC1.txt")
	if error != nil {
		fmt.Print(error)
	}
	defer file.Close()
	content, error := io.ReadAll(file)
	if error != nil {
		fmt.Print(error)
	}
	encodingString := base64.StdEncoding.EncodeToString(content)

	fmt.Print(encodingString)
}
