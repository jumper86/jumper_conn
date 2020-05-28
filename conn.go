package main

func NewtcpConnOptions(side int8, maxMsgSize int64, readTimeout int64, writeTimeout int64, asyncWriteSize int) *ConnOptions {
	return &ConnOptions{
		maxMsgSize:     maxMsgSize,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		asyncWriteSize: asyncWriteSize,
		side:           side,
	}
}

func NewwsConnOptions(side int8, maxMsgSize int64, readTimeout int64, writeTimeout int64, asyncWriteSize int,
	pongWait int64, pingPeriod int64, closeGracePeriod int64) *ConnOptions {
	return &ConnOptions{
		maxMsgSize:     maxMsgSize,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		asyncWriteSize: asyncWriteSize,
		side:           side,

		pongWait:         pongWait,
		pingPeriod:       pingPeriod,
		closeGracePeriod: closeGracePeriod,
	}
}
