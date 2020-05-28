package jumper_conn

import (
	"net"

	"github.com/gorilla/websocket"
	"github.com/jumper86/jumper_conn/impl/conn"
	"github.com/jumper86/jumper_conn/interf"
)

func NewwsConn(c *websocket.Conn, co interf.ConnOptionsInterf, handler interf.Handler) (interf.Conn, error) {
	wsConn, err := conn.CreatewsConn(c, co, handler)
	if err != nil {
		return nil, err
	}
	return wsConn, nil
}

func NewtcpConn(c net.Conn, co interf.ConnOptionsInterf, handler interf.Handler) (interf.Conn, error) {
	wsConn, err := conn.CreatetcpConn(c, co, handler)
	if err != nil {
		return nil, err
	}
	return wsConn, nil
}

func NewtcpConnOptions(side int8, maxMsgSize int64, readTimeout int64, writeTimeout int64, asyncWriteSize int64) interf.ConnOptionsInterf {
	option := conn.CreatetcpConnOptions(side, maxMsgSize, readTimeout, writeTimeout, asyncWriteSize)
	return option
}

func NewwsConnOptions(side int8, maxMsgSize int64, readTimeout int64, writeTimeout int64, asyncWriteSize int64,
	pongWait int64, pingPeriod int64, closeGracePeriod int64) interf.ConnOptionsInterf {
	option := conn.CreatewsConnOptions(side, maxMsgSize, readTimeout, writeTimeout, asyncWriteSize, pongWait, pingPeriod, closeGracePeriod)
	return option
}
