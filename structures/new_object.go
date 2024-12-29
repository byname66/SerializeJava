package structures

import (
	"fmt"
	"main/common"
)

type NewObject struct {
	TC_OBJECT  byte
	ClassDesc  *ClassDesc
	NewHandle  uint32
	ClassDatas *ClassDatas
}

func NewNewObject(classDesc *ClassDesc, newHandle uint32, classDatas *ClassDatas) *NewObject {
	return &NewObject{
		TC_OBJECT:  TC_OBJECT,
		ClassDesc:  classDesc,
		NewHandle:  newHandle,
		ClassDatas: classDatas,
	}
}

func ParseNewObject(parser *StructuresParser) (*NewObject, error) {
	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseNewObject:\n %v", err)
	}
	if signByte != TC_OBJECT {
		return nil, fmt.Errorf("in ParseNewObject: No TC_OBJECT")
	}
	classDesc, err := ParseClassDesc(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseNewObject:\n %v", err)
	}
	newHandle := parser.AddHandle()
	classDatas, err := ParseClassDatas(parser, classDesc)
	if err != nil {
		return nil, fmt.Errorf("in ParseNewObject:\n %v", err)
	}
	newObject := NewNewObject(classDesc, newHandle, classDatas)
	parser.AddReferenced(newObject)
	return newObject, nil
}

func (newObject *NewObject) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteByte(TC_OBJECT)
	err := newObject.ClassDesc.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteNewObject:\n %v", err)
	}
	err = newObject.ClassDatas.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteNewObject:\n %v", err)
	}
	return nil
}

func (newObject *NewObject) ToString(indent int) (string, error) {
	sb := common.NewStringBuilder(&indent)
	var result string
	result += sb.Build(" @NewObject")
	indent += IndentSpaceCount
	result += sb.Buildf("- TC_OBJECT  ", []interface{}{TC_OBJECT})
	str, err := newObject.ClassDesc.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in NewObject#ToString:\n%v", err)
	}
	result += str
	result += sb.Buildf("- NewHandle  ", []interface{}{newObject.NewHandle})
	str, err = newObject.ClassDatas.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in NewObject#ToString:\n%v", err)
	}
	result += str
	return result, nil
}

func (newObject *NewObject) GetNewHandle() uint32 {
	return newObject.NewHandle
}

type ClassDatas struct {
	AllDescs   []*ClassDesc
	ClassDatas []*ClassData
}

// func ParseClassDatas(parser *StructuresParser, classDesc *ClassDesc) (*ClassDatas, error) {
// 	var classDatas []*ClassData

// 	TclassDatas := new(ClassDatas)

// 	if classDesc.NewClassDesc == nil || classDesc.NewClassDesc.ClassDescInfo == nil {
// 		classData, err := ParseClassData4Proxy(parser)
// 		if err != nil {
// 			return nil, fmt.Errorf("in ParseNewObject:\n %v", err)
// 		}
// 		classDatas = append(classDatas, classData)
// 		TclassDatas.ClassDatas = classDatas
// 		return TclassDatas, nil
// 	} else {
// 		classData, err := ParseClassData(parser, classDesc.NewClassDesc.ClassDescInfo)
// 		if err != nil {
// 			return nil, fmt.Errorf("in ParseNewObject:\n %v", err)
// 		}
// 		classDatas = append(classDatas, classData)
// 		TclassDatas.ClassDatas = classDatas
// 		return TclassDatas, nil
// 	}
// }

// func ParseClassDatas(parser *StructuresParser, classDesc *ClassDesc) (*ClassDatas, error) {
// 	classDatas := new(ClassDatas)
// 	var cds []*ClassData
// 	allDescs, err := parser.FindAllDescs(classDesc)

// 	classDatas.AllDescs = allDescs
// 	if err != nil {
// 		return nil, fmt.Errorf("in ParseClassDatas:\n %v", err)
// 	}
// 	for i := len(allDescs) - 1; i >= 0; i-- {
// 		desc := allDescs[i]
// 		if desc.NewClassDesc.ClassDescInfo != nil {
// 			classData, err := ParseClassData(parser, desc.NewClassDesc.ClassDescInfo)
// 			if err != nil {
// 				return nil, fmt.Errorf("in ParseClassDatas:\n %v", err)
// 			}
// 			cds = append(cds, classData)
// 		} else if desc.NewClassDesc.ProxyClassDescInfo != nil {
// 			classData, err := ParseClassData4Proxy(parser)
// 			if err != nil {
// 				return nil, fmt.Errorf("in ParseClassDatas:\n %v", err)
// 			}
// 			cds = append(cds, classData)
// 		} else {
// 			return nil, fmt.Errorf("in ParseClassDatas:The ClassDesc is NULL")
// 		}

// 	}
// 	classDatas.ClassDatas = cds
// 	return classDatas, nil
// }

func ParseClassDatas(parser *StructuresParser, classDesc *ClassDesc) (*ClassDatas, error) {
	classDatas := new(ClassDatas)
	var cds []*ClassData
	allDescs, err := parser.FindAllDescs(classDesc)

	classDatas.AllDescs = allDescs
	if err != nil {
		return nil, fmt.Errorf("in ParseClassDatas:\n %v", err)
	}
	for i := len(allDescs) - 1; i >= 0; i-- {
		desc := allDescs[i]
		if desc.NewClassDesc.ClassDescInfo != nil {
			classData, err := ParseClassData(parser, desc.NewClassDesc.ClassDescInfo)
			if err != nil {
				return nil, fmt.Errorf("in ParseClassDatas:\n %v", err)
			}
			cds = append(cds, classData)
		} else {
			return nil, fmt.Errorf("in ParseClassDatas:The ClassDesc is NULL")
		}

	}
	classDatas.ClassDatas = cds
	return classDatas, nil
}

func (cds *ClassDatas) ToByte(parser *StructuresParser) error {
	for i := 0; i < len(cds.ClassDatas); i++ {
		err := cds.ClassDatas[i].ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteClassDatas:\n %v", err)
		}
	}
	return nil
}

//	func (cds *ClassDatas) ToString(indent int) (string, error) {
//		var (
//			result string
//			str    string
//			err    error
//		)
//		sb := common.NewStringBuilder(&indent)
//		result += sb.Build(" @ClassDatas : []ClassData")
//		indent += IndexToArraySpaceCount
//		//AllDescs' every desc is about "NewClassDesc" so they all have the attribute `className`(selected in structures_parser#FindAllDescs)
//		for i, j := len(cds.AllDescs)-1, 0; i >= 0 && j < len(cds.AllDescs); i--,j++{
//			result += sb.Buildf("- className  ", []interface{}{cds.AllDescs[j].NewClassDesc.ClassName})
//			str, err = cds.ClassDatas[i].ToString(indent+IndexToArraySpaceCount, cds.AllDescs[j])
//			if err != nil {
//				return "", fmt.Errorf("in ClassDatas#ToString\n%v", err)
//			}
//			result += str
//		}
//		return result, err
//	}
func (cds *ClassDatas) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)

	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @ClassDatas : []ClassData")
	indent += IndexToArraySpaceCount
	if len(cds.AllDescs) != len(cds.ClassDatas) {
		return "", fmt.Errorf("length mismatch: AllDescs(%d) vs ClassDatas(%d)", len(cds.AllDescs), len(cds.ClassDatas))
	}
	for i, j := len(cds.AllDescs)-1, 0; i >= 0 && j < len(cds.AllDescs); i, j = i-1, j+1 {

		result += sb.Buildf("- className  ", []interface{}{cds.AllDescs[i].NewClassDesc.ClassName})
		str, err = cds.ClassDatas[j].ToString(indent+IndexToArraySpaceCount, cds.AllDescs[i])
		if err != nil {
			return "", fmt.Errorf("in ClassDatas#ToString\n%v", err)
		}

		result += str
	}

	return result, nil
}
