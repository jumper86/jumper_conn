package cst

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
