package conn

import (
	"github.com/jumper86/jumper_conn/def"
	"github.com/jumper86/jumper_conn/interf"
)

type ConnOptions struct {
	maxMsgSize     int64
	readTimeout    int64
	writeTimeout   int64
	asyncWriteSize int64
	side           int8

	pongWait         int64
	pingPeriod       int64
	closeGracePeriod int64
}

func (this *ConnOptions) GetMaxMsgSize() int64 {
	return this.maxMsgSize
}
func (this *ConnOptions) GetReadTimeout() int64 {
	return this.readTimeout
}
func (this *ConnOptions) GetWriteTimeout() int64 {
	return this.writeTimeout
}
func (this *ConnOptions) GetAsyncWriteSize() int64 {
	return this.asyncWriteSize
}
func (this *ConnOptions) GetSide() int8 {
	return this.side
}
func (this *ConnOptions) GetPongWait() int64 {
	return this.pongWait
}
func (this *ConnOptions) GetPingPeriod() int64 {
	return this.pingPeriod
}
func (this *ConnOptions) GetCloseGracePeriod() int64 {
	return this.closeGracePeriod
}

func (this *ConnOptions) CheckValid() error {
	if this.asyncWriteSize <= 0 {
		return def.ErrInvalidConnParam
	}
	if this.pingPeriod != 0 && this.pingPeriod >= this.pongWait {
		return def.ErrInvalidConnParam
	}
	return nil
}

func CreatetcpConnOptions(side int8, maxMsgSize int64, readTimeout int64, writeTimeout int64, asyncWriteSize int64) interf.ConnOptionsInterf {
	return &ConnOptions{
		maxMsgSize:     maxMsgSize,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		asyncWriteSize: asyncWriteSize,
		side:           side,
	}
}

func CreatewsConnOptions(side int8, maxMsgSize int64, readTimeout int64, writeTimeout int64, asyncWriteSize int64,
	pongWait int64, pingPeriod int64, closeGracePeriod int64) interf.ConnOptionsInterf {
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
