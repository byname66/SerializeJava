package structures

import (
	"fmt"
	"main/common"
)

type Annotation struct {
	Flag         byte
	EndBlockData byte
	Contents     []*Content
}

func ParseAnnotation(parser *StructuresParser) (*Annotation, error) {
	peekByte, err := parser.ByteReader.PeekByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseAnnotation:\n %w", err)
	}
	Annotation := new(Annotation)
	Annotation.Flag = peekByte
	var contents []*Content
	switch peekByte {
	case TC_ENDBLOCKDATA:

		Annotation.EndBlockData = TC_ENDBLOCKDATA
		err := parser.ByteReader.JumpByte()
		if err != nil {
			return nil, fmt.Errorf("in ParseAnnotation:\n %w", err)
		}

		//If Annotation contains content:
	case TC_BLOCKDATA, TC_BLOCKDATALONG, TC_OBJECT, TC_CLASS, TC_ARRAY, TC_STRING, TC_LONGSTRING, TC_ENUM, TC_CLASSDESC, TC_REFERENCE, TC_NULL, TC_EXCEPTION, TC_RESET:
		for {
			content, err := ParseContent(parser)
			if err != nil {
				return nil, fmt.Errorf("in ParseAnnotation:\n %w", err)
			}
			contents = append(contents, content)
			peekByte, err := parser.ByteReader.PeekByte()
			if err != nil {
				return nil, fmt.Errorf("in ParseAnnotation:\n %w", err)
			}
			//Read the TC_ENDBLOCKDATA,then finish parse Content.
			if peekByte == TC_ENDBLOCKDATA {
				err = parser.ByteReader.JumpByte()
				if err != nil {
					return nil, fmt.Errorf("in ParseAnnotation:\n %w", err)
				}
				Annotation.EndBlockData = TC_ENDBLOCKDATA
				break
			}
		}
		Annotation.Contents = contents
	}

	return Annotation, nil
}

func (an *Annotation) ToByte(parser *StructuresParser) error {
	if an.Contents != nil {
		for i := 0; i < len(an.Contents); i++ {
			an.Contents[i].ToByte(parser)
		}
		parser.ByteReader.WriteByte(TC_ENDBLOCKDATA)
		return nil
	} else if an.EndBlockData != 0 {
		parser.ByteReader.WriteByte(TC_ENDBLOCKDATA)
		return nil
	} else {
		return fmt.Errorf("in WriteAnnotation:No matched field")
	}
}

func (an *Annotation) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build("@Annotation")
	indent += IndentSpaceCount
	switch an.Flag {
	case TC_ENDBLOCKDATA:
		result += sb.Buildf("- TC_ENDBLOCKDATA  ", []interface{}{TC_ENDBLOCKDATA})
	default:
		result += sb.Build("@Contents")
		indent += IndexToArraySpaceCount
		for i := 0; i < len(an.Contents); i++ {
			result += sb.Buildf("Index  [", []interface{}{i, "]"})
			str, err = an.Contents[i].ToString(indent + IndentSpaceCount)
			if err != nil {
				return "", fmt.Errorf("in Annotation#ToString:\n%v", err)
			}
			result += str
		}
		indent -= IndexToArraySpaceCount
		result += sb.Buildf("- TC_ENDBLOCKDATA  ", []interface{}{TC_ENDBLOCKDATA})
	}
	return result, nil
}
