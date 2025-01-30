package chunk

type BinaryChunk struct {
	header                  // 头部
	sizeUpValues byte       // 主函数upValue数量
	mainFunc     *ProtoType // 主函数原型
}

type UpValue struct {
	InStack byte
	Idx     byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSZIET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

type header struct {
	signature       [4]byte // 签名
	version         byte    // 版本号 大版本号 x 16 + 小版本号
	format          byte    // chunk 格式号 默认为 0
	luacData        [6]byte // 校验
	cIntSize        byte    // c int 字节数
	sizeTSize       byte    // size t 字节数
	instructionSize byte    // lua 虚拟机指令字节数
	luaIntegerSize  byte    // lua int 字节数
	luaNumberSize   byte    // lua number 字节数
	luacInt         int64   // 检测大小端
	luacNum         float64 // 检测浮点数格式
}

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

type ProtoType struct {
	Source          string // 源文件名字 以 @ 开头则为 Lua 源文件编译来的，以 = 开头则为标准输入编译来的
	LineDefined     uint32 // 起止行号
	LastLineDefined uint32
	NumParams       byte          // 固定参数个数
	IsVararg        byte          // 是否有变长参数
	MaxStackSize    byte          // 寄存器数量 但是 lua 实际使用栈结构来模拟寄存器
	Code            []uint32      // 指令表
	Constants       []interface{} // 常量表
	UpValues        []UpValue     // up value 表
	ProtoTypes      []*ProtoType  // 子函数原型
	LineInfo        []uint32      // 行号表
	LocVars         []LocVar      // 局部变量表
	UpValueNames    []string      // up value 名称表格
}

func UnDump(data []byte) *ProtoType {
	reader := &reader{data}
	reader.checkHeader()        // 校验头部
	reader.readByte()           // 跳过UpValue数量
	return reader.readProto("") // 读取函数原型
}
