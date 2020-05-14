package impl

import (
	"github.com/jumper86/jumper_conn/cst"
	"github.com/jumper86/jumper_conn/interf"
)

type Side int8

const (
	ServerSide Side = iota
	ClientSide
)

const (
	maxMsgSize     = 8192
	readTimeout    = 10
	writeTimeout   = 10
	asyncWriteSize = 20

	pongWait         = 60
	pingPeriod       = (pongWait * 9) / 10
	closeGracePeriod = 10
)

type connOptions struct {
	maxMsgSize     int64
	readTimeout    int64
	writeTimeout   int64
	asyncWriteSize int
	side           int8

	pongWait         int64
	pingPeriod       int64
	closeGracePeriod int64
}

func newTcpConnOptions(side int8, maxMsgSize int64, readTimeout int64, writeTimeout int64, asyncWriteSize int) *connOptions {
	return &connOptions{
		maxMsgSize:     maxMsgSize,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		asyncWriteSize: asyncWriteSize,
		side:           side,
	}
}

func newWsConnOptions(side int8, maxMsgSize int64, readTimeout int64, writeTimeout int64, asyncWriteSize int,
	pongWait int64, pingPeriod int64, closeGracePeriod int64) *connOptions {
	return &connOptions{
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

func checkOp(co *connOptions, handler interf.Handler) error {
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
