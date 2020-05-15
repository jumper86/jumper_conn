package interf

//定义打包和加密类型
type OperationType int

const (
	PackageOpMin OperationType = 0 + iota
	//封包
	PacketBase64
	PacketJson
	PacketXml
	PacketProtobuf

	//压缩
	CompressGzip
	CompressZlib

	//加密
	EncryptMd5
	EncryptSha1
	EncryptAes
	EncryptDes
	EncryptRsa

	PackageOpMax
)

//操作接口
type Operation interface {
	//Init(direct bool, params []interface{}) bool //direct 表示操作方向, true　表示　编码/压缩/加密，　false 表示　解码/解压/解密
	Operate(input interface{}, output interface{}) (bool, error)
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

type PacketOp interface {
	Operation
	Pack(originData interface{}) ([]byte, error)
	Unpack(packData []byte, obj interface{}) error
}

type Transform interface {
	AddOp(opType OperationType, direct bool, params []interface{}) bool
	Execute(input interface{}, output interface{}) error
	Reset()
}
