package structures

import (
	"encoding/binary"
	"fmt"
	"main/common"
)

type BlockData struct {
	FLAG           byte
	Blockdatashort *BlockDataShort
	Blockdatalong  *BlockDataLong
}

func ParseBlockData(parser *StructuresParser) (*BlockData, error) {
	peekByte, err := parser.ByteReader.PeekByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
	}
	blockData := new(BlockData)
	switch peekByte {
	case TC_BLOCKDATA:
		blockdatashort, err := ParseBlockDataShort(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
		}
		blockData.Blockdatashort = blockdatashort
		blockData.FLAG = TC_BLOCKDATA
	case TC_BLOCKDATALONG:
		blockdatalong, err := ParseBlockDataLong(parser)
		if err != nil {
			return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
		}
		blockData.Blockdatalong = blockdatalong
		blockData.FLAG = TC_BLOCKDATALONG
	}
	return blockData, nil
}

func (bd *BlockData) ToByte(parser *StructuresParser) error {
	if bd.Blockdatashort != nil {
		err := bd.Blockdatashort.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteBlockData:\n %v", err)
		}
		return nil
	} else if bd.Blockdatalong != nil {
		err := bd.Blockdatalong.ToByte(parser)
		if err != nil {
			return fmt.Errorf("in WriteBlockData:\n %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("in WriteBlockData:No field in BlockData")
	}
}

func (bd *BlockData) ToString(indent int) (string, error) {
	var (
		result string
		str    string
		err    error
	)
	switch bd.FLAG {
	case TC_BLOCKDATA:
		str, err = bd.Blockdatashort.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in BlockData#ToString:\n%v", err)
		}
		result += str
	case TC_BLOCKDATALONG:
		str, err = bd.Blockdatalong.ToString(indent)
		if err != nil {
			return "", fmt.Errorf("in BlockData#ToString:\n%v", err)
		}
		result += str
	}
	return result, nil
}

type BlockDataShort struct {
	TC_BLOCKDATA byte
	Size         uint8
	Data         []byte
}

func NewBlockDataShort(size uint8, data []byte) *BlockDataShort {
	return &BlockDataShort{
		TC_BLOCKDATA: TC_BLOCKDATA,
		Size:         size,
		Data:         data,
	}
}

func ParseBlockDataShort(parser *StructuresParser) (*BlockDataShort, error) {
	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
	}
	if signByte != TC_BLOCKDATA {
		return nil, fmt.Errorf("in ParseBlockDataShort: No TC_BLOCKDATA")
	}
	uint8Data, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
	}
	size := uint8(uint8Data)
	data, err := parser.ByteReader.ReadNByte(int(size))
	if err != nil {
		return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
	}
	return NewBlockDataShort(size, data), nil
}

func (bds *BlockDataShort) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteByte(TC_BLOCKDATA)
	parser.ByteReader.WriteNumber(bds.Size)
	parser.ByteReader.WriteNByte(bds.Data)
	return nil
}

func (bds *BlockDataShort) ToString(indent int) (string, error) {
	var (
		result string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @BlockData")
	indent += IndentSpaceCount
	result += sb.Buildf("- TC_BLOCKDATA  ", []interface{}{TC_BLOCKDATA})
	byteArray, err := common.ConvertNumberToBytes(bds.Size)
	if err != nil {
		return "", fmt.Errorf("in BlockDataShort#ToString:\n%v", err)
	}
	result += sb.Buildf("- size  ", []interface{}{bds.Size, "  -  ", byteArray})
	result += sb.Buildf("- data  ", []interface{}{bds.Data})
	return result, nil
}

type BlockDataLong struct {
	TC_blockdatalong byte
	Size             int32
	Data             []byte
}

func NewBlockDataLong(size int32, data []byte) *BlockDataLong {
	return &BlockDataLong{
		TC_blockdatalong: TC_BLOCKDATALONG,
		Size:             size,
		Data:             data,
	}
}

func ParseBlockDataLong(parser *StructuresParser) (*BlockDataLong, error) {
	signByte, err := parser.ByteReader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
	}
	if signByte != TC_BLOCKDATA {
		return nil, fmt.Errorf("in ParseBlockDataShort: No TC_BLOCKDATA")
	}
	byteArray, err := parser.ByteReader.ReadNByte(4)
	if err != nil {
		return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
	}
	size := int32(binary.BigEndian.Uint32(byteArray))
	data, err := parser.ByteReader.ReadNByte(int(size))
	if err != nil {
		return nil, fmt.Errorf("in ParseBlockData:\n %v", err)
	}
	return NewBlockDataLong(size, data), nil
}

func (bdl *BlockDataLong) ToByte(parser *StructuresParser) error {
	parser.ByteReader.WriteByte(TC_BLOCKDATALONG)
	parser.ByteReader.WriteNumber(bdl.Size)
	parser.ByteReader.WriteNByte(bdl.Data)
	return nil
}

func (bds *BlockDataLong) ToString(indent int) (string, error) {
	var (
		result string
		err    error
	)
	sb := common.NewStringBuilder(&indent)
	result += sb.Build(" @BlockDataLong")
	indent += IndentSpaceCount
	result += sb.Buildf("- TC_BLOCKDATALONG  ", []interface{}{TC_BLOCKDATALONG})
	byteArray, err := common.ConvertNumberToBytes(bds.Size)
	if err != nil {
		return "", fmt.Errorf("in BlockData#ToString:\n%v", err)
	}
	result += sb.Buildf("- size  ", []interface{}{bds.Size, "  -  ", byteArray})
	result += sb.Buildf("- data  ", []interface{}{bds.Data})
	return result, nil
}
