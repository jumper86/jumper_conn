package interf

import (
	"net"
)

type Conn interface {
	GetConn() net.Conn
	Close()
	IsClosed() bool

	Write(data []byte) error
	AsyncWrite(data []byte) error

	LocalAddr() net.Addr
	RemoteAddr() net.Addr

	Set(string, interface{})
	Get(string) interface{}
	Del(string)

	Run()
}

type ConnOptionsInterf interface {
	GetMaxMsgSize() int64
	GetReadTimeout() int64
	GetWriteTimeout() int64
	GetAsyncWriteSize() int64
	GetSide() int8
	GetPongWait() int64
	GetPingPeriod() int64
	GetCloseGracePeriod() int64

	CheckValid() error
}
