package conn

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jumper86/jumper_conn/interf"

	"github.com/gorilla/websocket"
	"github.com/jumper86/jumper_conn/def"
)

type wsConn struct {
	closed      int32
	writeBuffer chan []byte
	closeChan   chan struct{}
	ctx         map[string]interface{}

	conn    *websocket.Conn
	co      interf.ConnOptionsInterf
	handler interf.Handler
}

func CreatewsConn(conn *websocket.Conn, co interf.ConnOptionsInterf, handler interf.Handler) (interf.Conn, error) {

	rc := &wsConn{
		conn:        conn,
		closed:      0,
		writeBuffer: make(chan []byte, co.GetAsyncWriteSize()),
		closeChan:   make(chan struct{}),
		co:          co,
		ctx:         make(map[string]interface{}),
		handler:     handler,
	}

	return rc, nil
}

func (this *wsConn) Run() {
	this.run()
}

func (this *wsConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *wsConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *wsConn) GetConn() net.Conn {
	return this.conn.UnderlyingConn()
}

func (this *wsConn) Close() {

	this.close(nil)
}
func (this *wsConn) IsClosed() bool {
	return atomic.LoadInt32(&this.closed) == 1
}

func (this *wsConn) Write(data []byte) error {
	closed := this.IsClosed()
	if closed {
		return def.ErrConnClosed
	}

	this.setWriteDeadline(this.co.GetWriteTimeout())
	defer this.setWriteDeadline(0)

	err := this.conn.WriteMessage(websocket.TextMessage, data)
	return err
}

func (this *wsConn) AsyncWrite(data []byte) (err error) {
	closed := this.IsClosed()
	if closed {
		return def.ErrConnClosed
	}

	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = def.ErrConnClosed
		}
	}()

	this.writeBuffer <- data
	return nil
}
func (this *wsConn) Set(key string, value interface{}) {
	this.ctx[key] = value
}

func (this *wsConn) Get(key string) interface{} {
	if value, ok := this.ctx[key]; ok {
		return value
	}
	return nil
}

func (this *wsConn) Del(key string) {
	delete(this.ctx, key)
}

////////////////////////////////////////////////////////////// impl
//服务端和客户端都需要
func (this *wsConn) setReadLimit() {
	this.conn.SetReadLimit(this.co.GetMaxMsgSize())
}

//服务端发送ping, 接收pong
//客户端接收ping, 发送pong, 默认底层处理已经使用回复了pong
func (this *wsConn) sendPing() {
	ticker := time.NewTicker(time.Duration(this.co.GetPingPeriod()))

	for {
		select {
		case <-this.closeChan:
			return
		case <-ticker.C:
			this.conn.WriteControl(websocket.PingMessage, nil,
				time.Now().Add(time.Duration(this.co.GetWriteTimeout())*time.Second))
		}
	}

}

func (this *wsConn) handlePong() {
	this.conn.SetPongHandler(func(appData string) error {
		return this.conn.SetReadDeadline(time.Now().Add(time.Duration(this.co.GetPongWait()) * time.Second))
	})
}

func (this *wsConn) setWriteDeadline(timeout int64) {
	if this.co.GetWriteTimeout() > 0 {
		this.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
}

func (this *wsConn) setReadDeadline(timeout int64) {
	if this.co.GetReadTimeout() > 0 {
		this.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
}

func (this *wsConn) close(err error) {
	if !atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		return
	}

	close(this.closeChan)
	close(this.writeBuffer)

	if err == nil || err == def.ErrConnClosed {
		content := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "byebye.")
		this.conn.WriteMessage(websocket.CloseMessage, content)
		time.Sleep(time.Duration(this.co.GetCloseGracePeriod()) * time.Second)
	}
	this.conn.Close()

	this.handler.OnClose(err)

	this.ctx = nil
	this.handler = nil

}

func (this *wsConn) asyncWrite(wg *sync.WaitGroup) error {

	wg.Done()

	var err error
readLoop:
	for {
		select {
		case <-this.closeChan:
			err = def.ErrConnClosed
			break readLoop

		case data, ok := <-this.writeBuffer:
			if !ok {
				err = def.ErrConnClosed
				break readLoop
			}

			this.setWriteDeadline(this.co.GetWriteTimeout())

			err = this.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					err = def.ErrConnClosed
				} else {
					err = def.ErrConnUnexpectedClosed
				}
				break readLoop
			}

			this.setWriteDeadline(0)
		}
	}

	this.close(err)
	return err
}

func (this *wsConn) read(wg *sync.WaitGroup) error {

	wg.Done()

	var err error

readLoop:
	for {
		select {
		case <-this.closeChan:
			err = def.ErrConnClosed
			break readLoop
		default:
			this.setReadDeadline(this.co.GetReadTimeout())
			_, msg, err := this.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					err = def.ErrConnClosed
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

func (this *wsConn) run() {

	if this.IsClosed() {
		return
	}

	this.setReadLimit()
	if this.co.GetSide() == def.ServerSide {
		go this.sendPing()
		this.handlePong()
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go this.read(wg)
	go this.asyncWrite(wg)
	wg.Done()
}
