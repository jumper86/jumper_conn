package impl

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/jumper86/jumper_conn/interf"

	"github.com/gorilla/websocket"
	"github.com/jumper86/jumper_conn/cst"
)

type WsConn struct {
	conn        *websocket.Conn
	closed      int32
	writeBuffer chan []byte
	closeChan   chan struct{}
	connOptions
	handler interf.Handler
}

func NewWsConn(conn *websocket.Conn, co connOptions, handler interf.Handler) (interf.JConn, error) {
	err := checkOp(co, handler)
	if err != nil {
		return nil, err
	}
	rc := &WsConn{
		conn:        conn,
		closed:      0,
		writeBuffer: make(chan []byte, co.asyncWriteSize),
		closeChan:   make(chan struct{}),
		connOptions: co,
		handler:     handler,
	}
	return rc, nil
}

func (this *WsConn) Run() {

	go this.read()
	go this.asyncWrite()

}

func (this *WsConn) GetConn() net.Conn {
	return this.conn.UnderlyingConn()
}

func (this *WsConn) Close() {

	this.close(nil)
}
func (this *WsConn) IsClosed() bool {
	closed := atomic.LoadInt32(&this.closed)
	return closed == 1
}

func (this *WsConn) close(err error) {
	if !atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		return
	}

	//todo: 释放资源
	close(this.closeChan)
	close(this.writeBuffer)

	if err == nil || err == cst.ErrConnClosed {
		content := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "byebye.")
		this.conn.WriteMessage(websocket.CloseMessage, content)
	}
	this.conn.Close()

	this.handler.OnClose(err)

}

func (this *WsConn) setWriteDeadline(timeout int64) {
	if this.writeTimeout > 0 {
		this.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
}

func (this *WsConn) setReadDeadline(timeout int64) {
	if this.readTimeout > 0 {
		this.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
}

func (this *WsConn) Write(data []byte) error {
	closed := this.IsClosed()
	if closed {
		return cst.ErrConnClosed
	}

	this.setWriteDeadline(this.writeTimeout)
	defer this.setWriteDeadline(0)

	err := this.conn.WriteMessage(websocket.TextMessage, data)
	return err
}

func (this *WsConn) AsyncWrite(data []byte) (err error) {
	closed := this.IsClosed()
	if closed {
		return cst.ErrConnClosed
	}

	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = cst.ErrConnClosed
		}
	}()

	this.writeBuffer <- data
	return nil
}

func (this *WsConn) asyncWrite() error {

	var err error
readLoop:
	for {
		select {
		case <-this.closeChan:
			err = cst.ErrConnClosed
			break readLoop

		case data, ok := <-this.writeBuffer:
			if !ok {
				err = cst.ErrConnClosed
				break readLoop
			}

			this.setWriteDeadline(this.writeTimeout)

			err = this.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					err = cst.ErrConnClosed
				} else {
					err = cst.ErrConnUnexpectedClosed
				}
				break readLoop
			}

			this.setWriteDeadline(0)
		}
	}

	this.close(err)
	return err
}

func (this *WsConn) read() error {

	var err error

readLoop:
	for {
		select {
		case <-this.closeChan:
			err = cst.ErrConnClosed
			break readLoop
		default:
			this.setReadDeadline(this.readTimeout)
			_, msg, err := this.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					err = cst.ErrConnClosed
				}
				break readLoop
			}
			this.setReadDeadline(0)

			err = this.handler.OnMessage(msg)
			if err != nil {
				break readLoop
			}
		}

	}

	this.close(err)
	return err
}
