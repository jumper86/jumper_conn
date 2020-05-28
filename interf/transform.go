package interf

//定义打包和加密类型

//操作接口
type Operation interface {
	//Init(direct bool, params []interface{}) bool //direct 表示操作方向, true　表示　编码/压缩/加密，　false 表示　解码/解压/解密
	Operate(direct int8, input interface{}, output interface{}) (bool, error)
}

type CompressOp interface {
	Operation
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

type EncryptOp interface {
	Operation
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type Message struct {
	Type    uint16
	Content []byte
}

type PacketOp interface {
	Operation
	Pack(originData interface{}) ([]byte, error)
	Unpack(packData []byte, obj interface{}) error
}

type Transform interface {
	AddOp(opType int8, params []interface{}) bool
	Execute(direct int8, input interface{}, output interface{}) error
	Reset()
}
