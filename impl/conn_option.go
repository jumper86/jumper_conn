package impl

import (
	"github.com/jumper86/jumper_conn/cst"
	"github.com/jumper86/jumper_conn/interf"
)

type connOptions struct {
	readTimeout    int64
	writeTimeout   int64
	asyncWriteSize int
}

func checkOp(co connOptions, handler interf.Handler) error {
	if co.asyncWriteSize <= 0 {
		return cst.ErrCreateConnInvalidParam
	}
	if handler == nil {
		return cst.ErrCreateConnInvalidParam
	}
	return nil
}
