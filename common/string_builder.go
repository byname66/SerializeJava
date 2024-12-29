package common

import (
	"fmt"
	"strings"
)

type StringBuilder struct {
	indent *int
}

func NewStringBuilder(indent *int) *StringBuilder {
	return &StringBuilder{
		indent: indent,
	}
}

func (sb *StringBuilder) Buildf(str string, values []interface{}) string {
	var builder strings.Builder
	builder.WriteString(str)

	for _, value := range values {
		switch v := value.(type) {
		case byte:
			builder.WriteString(fmt.Sprintf("0x%02x", v))
		case []byte:
			builder.WriteString("0x")
			for i, b := range v {
				if i >= 1 {
					builder.WriteString(" ")
				}
				builder.WriteString(fmt.Sprintf("%02x", b))
			}
		default:
			builder.WriteString(fmt.Sprintf("%v", v))
		}
	}

	return sb.Build(builder.String())
}
func (sb *StringBuilder) Build(str string) string {
	spaces := strings.Repeat(" ", *sb.indent)
	str = spaces + str + "\n"
	return str
}

func (sb *StringBuilder) BuildWithSpaces(str string, count int) string {
	spaces := strings.Repeat(" ", count)
	str = spaces + str + "\n"
	return str
}
