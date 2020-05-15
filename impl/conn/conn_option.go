package conn

import (
	"github.com/jumper86/jumper_conn/cst"
	"github.com/jumper86/jumper_conn/interf"
)

type ConnOptions struct {
	maxMsgSize     int64
	readTimeout    int64
	writeTimeout   int64
	asyncWriteSize int
	side           int8

	pongWait         int64
	pingPeriod       int64
	closeGracePeriod int64
}

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

func checkOp(co *ConnOptions, handler interf.Handler) error {
	if co.asyncWriteSize <= 0 {
		return cst.ErrCreateConnInvalidParam
	}
	if handler == nil {
		return cst.ErrCreateConnInvalidParam
	}
	if co.pingPeriod != 0 && co.pingPeriod >= co.pongWait {
		return cst.ErrCreateConnInvalidParam
	}
	return nil
}
