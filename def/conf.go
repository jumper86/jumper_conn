package def

const (
	MaxMsgSize     = 8192
	ReadTimeout    = 600
	WriteTimeout   = 600
	AsyncWriteSize = 20

	PongWait         = 60
	PingPeriod       = (PongWait * 9) / 10
	CloseGracePeriod = 1
)
