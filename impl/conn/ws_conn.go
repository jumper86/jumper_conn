package conn

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jumper86/jumper_conn/interf"

	"github.com/gorilla/websocket"
	"github.com/jumper86/jumper_conn/cst"
)

type WsConn struct {
	connOptions
	conn        *websocket.Conn
	closed      int32
	writeBuffer chan []byte
	closeChan   chan struct{}

	handler interf.Handler
	ctx     map[string]interface{}
}

func NewWsConn(conn *websocket.Conn, co *connOptions, handler interf.Handler) (interf.Conn, error) {
	err := checkOp(co, handler)
	if err != nil {
		return nil, err
	}
	rc := &WsConn{
		conn:        conn,
		closed:      0,
		writeBuffer: make(chan []byte, co.asyncWriteSize),
		closeChan:   make(chan struct{}),
		connOptions: *co,
		handler:     handler,
	}

	rc.run()
	return rc, nil
}

func (this *WsConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}
func (this *WsConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *WsConn) GetConn() net.Conn {
	return this.conn.UnderlyingConn()
}

func (this *WsConn) Close() {

	this.close(nil)
}
func (this *WsConn) IsClosed() bool {
	return atomic.LoadInt32(&this.closed) == 1
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
func (this *WsConn) Set(key string, value interface{}) {
	this.ctx[key] = value
}

func (this *WsConn) Get(key string) interface{} {
	if value, ok := this.ctx[key]; ok {
		return value
	}
	return nil
}

func (this *WsConn) Del(key string) {
	delete(this.ctx, key)
}

////////////////////////////////////////////////////////////// impl
//服务端和客户端都需要
func (this *WsConn) setReadLimit() {
	this.conn.SetReadLimit(this.maxMsgSize)
}

//服务端发送ping, 接收pong
//客户端接收ping, 发送pong, 默认底层处理已经使用回复了pong
func (this *WsConn) sendPing() {
	ticker := time.NewTicker(time.Duration(this.pingPeriod))
	for _ = range ticker.C {
		this.conn.WriteControl(websocket.PingMessage, nil,
			time.Now().Add(time.Duration(this.writeTimeout)*time.Second))
	}

}

func (this *WsConn) handlePong() {
	this.conn.SetPongHandler(func(appData string) error {
		return this.conn.SetReadDeadline(time.Now().Add(time.Duration(this.pongWait) * time.Second))
	})
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

func (this *WsConn) asyncWrite(wg sync.WaitGroup) error {

	wg.Done()

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

func (this *WsConn) read(wg sync.WaitGroup) error {

	wg.Done()

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

func (this *WsConn) run() {

	if this.IsClosed() {
		return
	}

	this.setReadLimit()
	if this.side == ServerSide {
		this.sendPing()
		this.handlePong()
	}

	wg := sync.WaitGroup{}
	go this.read(wg)
	go this.asyncWrite(wg)
	wg.Done()
}
