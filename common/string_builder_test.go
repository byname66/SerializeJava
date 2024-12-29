package common

import (
	"fmt"
	"testing"
)

// import (
// 	"fmt"
// 	"main/structures"
// 	"testing"
// )

func TestBuildf(t *testing.T) {
	TC_OBJECT := byte(65)
	indent := 0
	sb := NewStringBuilder(&indent)
	var result string
	result += sb.Build(" @NewObject")
	indent += 4
	result += sb.Buildf("- TC_OBJECT  ", []interface{}{TC_OBJECT})
	fmt.Print(result)
}
