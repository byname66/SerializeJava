package structures

import (
	"fmt"
	"main/common"
)

type ProxyClassDescInfo struct {
	Count               int32
	ProxyInterfaceNames []*UTF
	ClassAnnotation     *Annotation
	SuperClassDesc      *ClassDesc
}

func NewProxyClassDescInfo(count int32, proxyInterfaceNames []*UTF, classAnnotation *Annotation, superClassDesc *ClassDesc) *ProxyClassDescInfo {
	return &ProxyClassDescInfo{
		Count:               count,
		ProxyInterfaceNames: proxyInterfaceNames,
		ClassAnnotation:     classAnnotation,
		SuperClassDesc:      superClassDesc,
	}
}

func ParseProxyClassDescInfo(parser *StructuresParser) (*ProxyClassDescInfo, error) {
	count, err := parser.ByteReader.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("in ParseProxyClassDescInfo:\n %v", err)
	}
	var proxyInterfaceNames []*UTF
	for i := 0; i < int(count); i++ {
		proxyInterfaceName, err := ParseProxyInterfaceName(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseProxyClassDescInfo:\n %v", err)
		}
		proxyInterfaceNames = append(proxyInterfaceNames, proxyInterfaceName)
	}
	classAnnotation, err := ParseAnnotation(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseProxyClassDescInfo:\n %v", err)
	}
	superClassDesc, err := ParseClassDesc(parser)
	if err != nil {
		return nil, fmt.Errorf("in ParseProxyClassDescInfo:\n %v", err)
	}
	return NewProxyClassDescInfo(count, proxyInterfaceNames, classAnnotation, superClassDesc), nil
}

func (pdi *ProxyClassDescInfo) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteNumber(pdi.Count)
	for i := 0; i < len(pdi.ProxyInterfaceNames); i++ {
		err := pdi.ProxyInterfaceNames[i].ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteProxyClassDescInfo:\n %v", err)
		}
	}
	err := pdi.ClassAnnotation.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteProxyClassDescInfo:\n %v", err)
	}
	err = pdi.SuperClassDesc.ToByte(parser)
	if err != nil {
		return fmt.Errorf("in WriteProxyClassDescInfo:\n %v", err)
	}
	return nil
}

func (pdi *ProxyClassDescInfo) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	byteArray, err := common.ConvertNumberToBytes(pdi.Count)
	if err != nil {
		return "", fmt.Errorf("in ProxyClassDescInfo#ToString:\n%v", err)
	}
	result += sb.Buildf("- count  ", []interface{}{pdi.Count, "  -  ", byteArray})
	result += sb.Build(" @ProxyInterfaceNames(Ut)")
	indent += IndexToArraySpaceCount
	for i := 0; i < len(pdi.ProxyInterfaceNames); i++ {
		result += sb.Buildf("Index  [", []interface{}{i, "]"})
		str, err = pdi.ProxyInterfaceNames[i].ToString(indent + IndentSpaceCount)
		if err != nil {
			return "", fmt.Errorf("in ProxyClassDescInfo#ToString:\n%v", err)
		}
		result += str
	}
	indent -= IndexToArraySpaceCount
	str, err = pdi.ClassAnnotation.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in ProxyClassDescInfo#ToString:\n%v", err)
	}
	result += str
	str, err = pdi.SuperClassDesc.ToString(indent)
	if err != nil {
		return "", fmt.Errorf("in ProxyClassDescInfo#ToString:\n%v", err)
	}
	result += str
	return result, nil
}

func ParseProxyInterfaceName(parser *StructuresParser) (*UTF, error) {
	length, err := parser.ByteReader.ReadInt16()
	if err != nil {
		return nil, fmt.Errorf(" in ParseProxyInterfaceName:\n %v", err)
	}
	name, byteArray, err := parser.utfReader.ReadUtf(int(length), &parser.ByteReader)
	if err != nil {
		return nil, fmt.Errorf(" in ParseProxyInterfaceName:\n %v", err)
	}
	return &UTF{
		Length:    length,
		Value:     name,
		ByteArray: byteArray,
	}, nil
}
