package conn

import (
	"github.com/jumper86/jumper_conn/def"
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

func checkOp(co *ConnOptions, handler interf.Handler) error {
	if co.asyncWriteSize <= 0 {
		return def.ErrInvalidConnParam
	}
	if handler == nil {
		return def.ErrInvalidConnParam
	}
	if co.pingPeriod != 0 && co.pingPeriod >= co.pongWait {
		return def.ErrInvalidConnParam
	}
	return nil
}
