package def

const (
	ServerSide int8 = iota
	ClientSide
)

const (
	MaxMsgSize     = 8192
	ReadTimeout    = 10
	WriteTimeout   = 10
	AsyncWriteSize = 20

	PongWait         = 60
	PingPeriod       = (PongWait * 9) / 10
	CloseGracePeriod = 1
)

const (
	TcpHeadSize = 4
)

const (
	PackageOpMin int8 = 0 + iota
	//封包
	PacketBase64
	PacketJson
	PacketXml
	PacketProtobuf
	PacketBinary

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

const (
	Forward  int8 = 1 //打包->压缩->加密
	Backward int8 = 2 //解密->解压->解包
)