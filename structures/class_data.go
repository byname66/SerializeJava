package structures

import (
	"fmt"
	"main/common"
)

type ClassData struct {
	NowrClass        *NowrClass
	WrClass          *WrClass
	ObjectAnnotation *Annotation
	IntFlag          int
	//externalContents *ExternalContents
}

// func ParseClassData(parser *StructuresParser, classDescInfo *ClassDescInfo) (*ClassData, error) {
// 	classData := new(ClassData)

//		if classDescInfo.HasFlag(SzaC_SERIALIZABLE) && !classDescInfo.HasFlag(SC_WRITE_METHOD) {
//			nowrClass, err := ParseNowrClass(parser, classDescInfo)
//			if err != nil {
//				return nil, fmt.Errorf("in ParseClassData:\n %v", err)
//			}
//			classData.NowrClass = nowrClass
//		} else if classDescInfo.HasFlag(SC_SERIALIZABLE) && classDescInfo.HasFlag(SC_WRITE_METHOD) {
//			wrClass, err := ParseWrClass(parser, classDescInfo)
//			if err != nil {
//				return nil, fmt.Errorf("in ParseClassData:\n %v", err)
//			}
//			classData.WrClass = wrClass
//			//Parse objectAnnotation
//			objectAnnotation, err := ParseAnnotation(parser)
//			if err != nil {
//				return nil, fmt.Errorf("in ParseClassData:\n %v", err)
//			}
//			classData.ObjectAnnotation = objectAnnotation
//		} else if classDescInfo.HasFlag(SC_EXTERNALIZABLE) && !classDescInfo.HasFlag(SC_BLOCK_DATA) {
//			return nil, fmt.Errorf("oooooooooooooooooooooo")
//		} else if classDescInfo.HasFlag(SC_EXTERNALIZABLE) && classDescInfo.HasFlag(SC_BLOCK_DATA) {
//			objectAnnotation, err := ParseAnnotation(parser)
//			if err != nil {
//				return nil, fmt.Errorf("in ParseClassData:\n %v", err)
//			}
//			classData.ObjectAnnotation = objectAnnotation
//		}
//		return classData, nil
//	}
func ParseClassData(parser *StructuresParser, classDescInfo *ClassDescInfo) (*ClassData, error) {
	classData := new(ClassData)

	if classDescInfo.HasFlag(SC_SERIALIZABLE) && !classDescInfo.HasFlag(SC_WRITE_METHOD) {
		nowrClass, err := ParseNowrClass(parser, classDescInfo)
		if err != nil {
			return nil, fmt.Errorf("in ParseClassData:\n %v", err)
		}
		classData.NowrClass = nowrClass
		classData.IntFlag = 1
	} else if classDescInfo.HasFlag(SC_SERIALIZABLE) && classDescInfo.HasFlag(SC_WRITE_METHOD) {
		wrClass, err := ParseWrClass(parser, classDescInfo)
		if err != nil {
			return nil, fmt.Errorf("in ParseClassData:\n %v", err)
		}
		classData.WrClass = wrClass
		//Parse objectAnnotation
		objectAnnotation, err := ParseAnnotation(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseClassData:\n %v", err)
		}
		classData.ObjectAnnotation = objectAnnotation
		classData.IntFlag = 2
	} else if classDescInfo.HasFlag(SC_EXTERNALIZABLE) && !classDescInfo.HasFlag(SC_BLOCK_DATA) {
		classData.IntFlag = 3
		return nil, fmt.Errorf("if happen this error,contact me")
	} else if classDescInfo.HasFlag(SC_EXTERNALIZABLE) && classDescInfo.HasFlag(SC_BLOCK_DATA) {
		objectAnnotation, err := ParseAnnotation(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseClassData:\n %v", err)
		}
		classData.ObjectAnnotation = objectAnnotation
		classData.IntFlag = 4
	}
	return classData, nil
}

// func ParseClassData4Proxy(parser *StructuresParser) (*ClassData, error) {
// 	classData := new(ClassData)
// 	nowrclass := new(NowrClass)
// 	var values []*Value
// 	value, err := ParseValue1(parser)
// 	if err != nil {
// 		return nil, fmt.Errorf("in ParseClassData:\n %v", err)
// 	}
// 	values = append(values, value)
// 	nowrclass.Values = values
// 	classData.NowrClass = nowrclass
// 	return classData, nil
// }

func (cd *ClassData) ToByte(parser *StructuresParser) error {
	if cd.NowrClass != nil {
		err := cd.NowrClass.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteClassData:\n %v", err)
		}
		return nil
	} else if cd.WrClass != nil {
		err := cd.WrClass.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteClassData:\n %v", err)
		}
		err = cd.ObjectAnnotation.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteClassData:\n %v", err)
		}
		return nil
	} else if cd.ObjectAnnotation != nil {
		err := cd.ObjectAnnotation.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteClassData:\n %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("in WriteClassData:no field in ClassData")
	}
}

func (cd *ClassData) ToString(indent int, classDesc *ClassDesc) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	switch cd.IntFlag {
	case 1:
		for i := 0; i < len(classDesc.NewClassDesc.ClassDescInfo.Fields.FieldDescs); i++ {
			f := classDesc.NewClassDesc.ClassDescInfo.Fields.FieldDescs[i]
			switch f.TypeCode {
			case "B", "C", "D", "F", "I", "J", "S", "Z":
				indent += IndentSpaceCount
				str, err = cd.NowrClass.Values[i].ToStringForClassData(indent)
				if err != nil {
					return "", fmt.Errorf("in ClassData#ToString:\n%v", err)
				}
				result += sb.Buildf("", []interface{}{"(", TypeTable[f.TypeCode], ")  ", f.PrimitiveDesc.FieldName.Value})
				result += str
				indent -= IndentSpaceCount
			case "[", "L":
				indent += IndentSpaceCount
				str, err = cd.NowrClass.Values[i].ToStringForClassData(indent)
				if err != nil {
					return "", fmt.Errorf("in ClassData#ToString:\n%v", err)
				}
				result += sb.Buildf("(", []interface{}{TypeTable[f.TypeCode], ")  ", f.ObjectDesc.FieldName.Value})
				result += sb.Buildf("", []interface{}{f.ObjectDesc.ClassName1.Value})

				result += str
				indent -= IndentSpaceCount
			}
		}
	case 2:
		for i := 0; i < len(classDesc.NewClassDesc.ClassDescInfo.Fields.FieldDescs); i++ {
			f := classDesc.NewClassDesc.ClassDescInfo.Fields.FieldDescs[i]
			switch f.TypeCode {
			case "B", "C", "D", "F", "I", "J", "S", "Z":
				indent += IndentSpaceCount
				str, err = cd.WrClass.NowrClass.Values[i].ToStringForClassData(indent)
				if err != nil {
					return "", fmt.Errorf("in ClassData#ToString:\n%v", err)
				}
				result += sb.Buildf("(", []interface{}{TypeTable[f.TypeCode], ")  ", f.PrimitiveDesc.FieldName.Value})
				result += str
				indent -= IndentSpaceCount
			case "[", "L":
				indent += IndentSpaceCount
				str, err = cd.WrClass.NowrClass.Values[i].ToStringForClassData(indent)
				if err != nil {
					return "", fmt.Errorf("in ClassData#ToString:\n%v", err)
				}
				result += sb.Buildf("(", []interface{}{TypeTable[f.TypeCode], ")  ", f.ObjectDesc.FieldName.Value})
				result += sb.Buildf("", []interface{}{f.ObjectDesc.ClassName1.Value})
				result += str
				indent -= IndentSpaceCount
			}
		}
		str, err = cd.ObjectAnnotation.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in ClassData#ToString:\n%v", err)
		}
		result += str
	case 3:
		return "", fmt.Errorf("not support EXTERNALIZABLE")
	case 4:
		str, err = cd.ObjectAnnotation.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in ClassData#ToString:\n%v", err)
		}
		result += str
	}
	return result, err
}

type NowrClass struct {
	Values []*Value
}

func ParseNowrClass(parser *StructuresParser, classDescInfo *ClassDescInfo) (*NowrClass, error) {
	fieldCount := classDescInfo.Fields.Count
	nowrClass := new(NowrClass)
	var (
		values []*Value
	)

	for i := 0; i < int(fieldCount); i++ {
		typeCode := (classDescInfo.Fields.FieldDescs)[i].TypeCode
		value, err := ParseValue(parser, typeCode)
		if err != nil {
			return nil, fmt.Errorf("in ParseNowrClass:\n %v", err)
		}
		values = append(values, value)
	}
	nowrClass.Values = values
	return nowrClass, nil
}

func (nc *NowrClass) ToByte(parser *StructuresParser) error {
	for i := 0; i < len(nc.Values); i++ {
		err := nc.Values[i].ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteNowrClass:\n %v", err)
		}
	}
	return nil
}

type WrClass struct {
	NowrClass *NowrClass
}

func ParseWrClass(parser *StructuresParser, classDescInfo *ClassDescInfo) (*WrClass, error) {

	nowrClass, err := ParseNowrClass(parser, classDescInfo)
	if err != nil {
		return nil, fmt.Errorf("in ParseWrClass:\n %v", err)
	}
	wrClass := new(WrClass)
	wrClass.NowrClass = nowrClass
	return wrClass, nil
}

func (wc *WrClass) ToByte(parser *StructuresParser) error {
	err := wc.NowrClass.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteWrClass:\n %v", err)
	}
	return nil
}
