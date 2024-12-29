package structures

var STREAM_MAGIC = []byte{0xAC, 0xED}

var STREAM_VERSION = []byte{0x00, 0x05}

// =============== TC_* ===============

// TC_BASE First tag value.
var TC_BASE byte = 0x70

// TC_NULL Null object reference.
var TC_NULL byte = 0x70

// TC_REFERENCE Reference to an object already written into the stream.
var TC_REFERENCE byte = 0x71

// TC_CLASSDESC new Class Descriptor.
var TC_CLASSDESC byte = 0x72

// TC_OBJECT new Object.
var TC_OBJECT byte = 0x73

// TC_STRING new String.
var TC_STRING byte = 0x74

// TC_ARRAY new Array.
var TC_ARRAY byte = 0x75

// TC_CLASS Reference to Class.
var TC_CLASS byte = 0x76

// TC_BLOCKDATA Block of optional data. Byte following tag indicates number of bytes in this block data.
var TC_BLOCKDATA byte = 0x77

// TC_ENDBLOCKDATA End of optional block data blocks for an object.
var TC_ENDBLOCKDATA byte = 0x78

// TC_RESET Reset stream context. All handles written into stream are reset.
var TC_RESET byte = 0x79

// TC_BLOCKDATALONG long Block data. The long following the tag indicates the number of bytes in this block data.
var TC_BLOCKDATALONG byte = 0x7A

// TC_EXCEPTION Exception during write.
var TC_EXCEPTION byte = 0x7B

// TC_LONGSTRING Long string.
var TC_LONGSTRING byte = 0x7C

// TC_PROXYCLASSDESC new Proxy Class Descriptor.
var TC_PROXYCLASSDESC byte = 0x7D

// TC_ENUM new Enum constant.
var TC_ENUM byte = 0x7E

// TC_MAX Last tag value.
var TC_MAX byte = 0x7F

// BASE_WRITE_HANDLE First wire handle to be assigned.
var BASE_WRITE_HANDLE uint32 = 0x7e0000

// =============== Bit Mask ===============

// SC_WRITE_METHOD Bit mask for ObjectStreamClass flag.
// Indicates a Serializable class defines its own writeObject method.
var SC_WRITE_METHOD byte = 0x01

// SC_SERIALIZABLE Bit mask for ObjectStreamClass flag. Indicates class is Serializable.
var SC_SERIALIZABLE byte = 0x02

// SC_EXTERNALIZABLE Bit mask for ObjectStreamClass flag. Indicates class is Externalizable.
var SC_EXTERNALIZABLE byte = 0x04

// SC_BLOCK_DATA Bit mask for ObjectStreamClass flag.
// Indicates Externalizable data written in Block Data mode. Added for PROTOCOL_VERSION_2.
var SC_BLOCK_DATA byte = 0x08

// SC_ENUM Bit mask for ObjectStreamClass flag. Indicates class is an enum type.
var SC_ENUM byte = 0x10

// self
var FLAG_CLASS_DESC byte = 0x02

var IndentSpaceCount int = 4
var IndexToArraySpaceCount int = 2

var SizeTable = map[string]int{
	"B": 1,
	"C": 2,
	"D": 8,
	"F": 4,
	"I": 4,
	"J": 8,
	"S": 2,
	"Z": 1,
}

var TypeTable = map[string]string{
	"B": "byte",
	"C": "char",
	"D": "double",
	"F": "float",
	"I": "integer",
	"J": "long",
	"S": "short",
	"Z": "boolean",
	"L": "Object",
	"[": "Array",
}
