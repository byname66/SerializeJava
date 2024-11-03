package structures

import "fmt"

type ClassAnnotation struct {
	endBlockData byte
	contents     []*Content
}

func ParseClassAnnotation(parser *StructuresParser) (*ClassAnnotation, error) {
	peekByte, err := parser.byteReader.PeekByte()
	if err != nil {
		return nil, fmt.Errorf("In ParseClassAnnotation:\n %w", err)
	}
	classAnnotation := new(ClassAnnotation)
	var contents []*Content
	switch peekByte {
	case TC_ENDBLOCKDATA:
		err := parser.byteReader.JumpByte()
		if err != nil {
			return nil, fmt.Errorf("In ParseClassAnnotation:\n %w", err)
		}
		classAnnotation.endBlockData = TC_ENDBLOCKDATA

		//If ClassAnnotation contains content:
	case TC_BLOCKDATA, TC_BLOCKDATALONG, TC_OBJECT, TC_CLASS, TC_ARRAY, TC_STRING, TC_LONGSTRING, TC_ENUM, TC_CLASSDESC, TC_REFERENCE, TC_NULL, TC_EXCEPTION, TC_RESET:
		content, err := ParseContent(parser)
		if err != nil {
			return nil, fmt.Errorf("In ParseClassAnnotation:\n %w", err)
		}
		contents = append(contents, content)
		signByte, err := parser.byteReader.PeekByte()
		if err != nil {
			return nil, fmt.Errorf("In ParseClassAnnotation:\n %w", err)
		}
		//Read the TC_ENDBLOCKDATA,then finish parse Content.
		if signByte == TC_ENDBLOCKDATA {
			err = parser.byteReader.JumpByte()
			if err != nil {
				return nil, fmt.Errorf("In ParseClassAnnotation:\n %w", err)
			}
			classAnnotation.endBlockData = TC_ENDBLOCKDATA
			break
		}
	}
	return classAnnotation, nil
}
