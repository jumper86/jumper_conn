package conn

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jumper86/jumper_conn/interf"

	"github.com/jumper86/jumper_conn/cst"
)

type tcpConn struct {
	ConnOptions
	conn        net.Conn
	closed      int32
	writeBuffer chan []byte
	closeChan   chan struct{}

	ctx     map[string]interface{}
	handler interf.Handler
}

func NewtcpConn(conn net.Conn, co *ConnOptions, handler interf.Handler) (interf.Conn, error) {
	err := checkOp(co, handler)
	if err != nil {
		return nil, err
	}
	rc := &tcpConn{
		conn:        conn,
		closed:      0,
		writeBuffer: make(chan []byte, co.asyncWriteSize),
		closeChan:   make(chan struct{}),
		ConnOptions: *co,
		ctx:         make(map[string]interface{}),
		handler:     handler,
	}

	return rc, nil
}

func (this *tcpConn) Run() {
	this.run()
}

func (this *tcpConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}
func (this *tcpConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *tcpConn) GetConn() net.Conn {
	return this.conn
}

func (this *tcpConn) Close() {
	this.close(nil)
}

func (this *tcpConn) IsClosed() bool {
	return atomic.LoadInt32(&this.closed) == 1
}

func (this *tcpConn) Write(data []byte) error {
	closed := this.IsClosed()
	if closed {
		return cst.ErrConnClosed
	}

	this.setWriteDeadline(this.writeTimeout)
	defer this.setWriteDeadline(0)

	length := len(data)
	written := 0
	var err error
	var l int

	for {
		l, err = this.conn.Write(data[written:])
		if err != nil {
			break
		}
		written += l
		if written == length {
			break
		}
	}

	return err
}

func (this *tcpConn) AsyncWrite(data []byte) (err error) {

	closed := this.IsClosed()
	if closed {
		return cst.ErrConnClosed
	}

	defer func() {
		if panicErr := recover(); panicErr != nil {
			//todo: 接收并且处理下面writeBuffer 写入之前已经关闭导致的panic
			err = cst.ErrConnClosed
			return
		}
	}()

	this.writeBuffer <- data
	return nil
}

func (this *tcpConn) Set(key string, value interface{}) {
	this.ctx[key] = value
}

func (this *tcpConn) Get(key string) interface{} {
	if value, ok := this.ctx[key]; ok {
		return value
	}
	return nil
}

func (this *tcpConn) Del(key string) {
	delete(this.ctx, key)
}

////////////////////////////////////////////////////////////// impl

func (this *tcpConn) setWriteDeadline(timeout int64) {
	if this.writeTimeout > 0 {
		this.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
}

func (this *tcpConn) setReadDeadline(timeout int64) {
	if this.readTimeout > 0 {
		this.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
}

func (this *tcpConn) close(err error) {
	swapped := atomic.CompareAndSwapInt32(&this.closed, 0, 1)
	if !swapped {
		return
	}
	//todo: clean resource
	close(this.closeChan)
	close(this.writeBuffer)

	this.conn.Close()

	this.handler.OnClose(err)

	this.ctx = nil
	this.handler = nil
}

func (this *tcpConn) asyncWrite(wg *sync.WaitGroup) error {

	wg.Done()

	var l int
	var err error

writeLoop:
	for {
		select {
		case <-this.closeChan:
			err = cst.ErrConnClosed
			break writeLoop
		case data, ok := <-this.writeBuffer:
			if !ok {
				err = cst.ErrConnClosed
				break writeLoop
			}

			length := len(data)
			written := 0

			for {

				this.setWriteDeadline(this.writeTimeout)
				l, err = this.conn.Write(data[written:])
				this.setWriteDeadline(0)

				if err != nil {
					break writeLoop
				}
				written += l

				if length == written {
					break
				}
			}
		}
	}

	this.close(err)
	return err
}

func (this *tcpConn) read(wg *sync.WaitGroup) (err error) {

	wg.Done()
readLoop:
	for {
		select {
		case <-this.closeChan:
			err = cst.ErrConnClosed
			break readLoop
		default:

			this.setReadDeadline(this.readTimeout)

			length := make([]byte, cst.TcpHeadSize)
			_, err = io.ReadFull(this.conn, length)
			if err != nil {
				break readLoop
			}
			left := binary.BigEndian.Uint32(length)

			content := make([]byte, left)
			_, err = io.ReadFull(this.conn, content)
			if err != nil {
				break readLoop
			}

			//process msg 可能会花较长时间，导致读超时断开
			this.setReadDeadline(0)

			err = this.handler.OnMessage(content)
			if err != nil {
				break readLoop
			}
		}
	}

	this.close(err)
	return
}

func (this *tcpConn) run() {
	if this.IsClosed() {
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := this.read(wg)
		if err != nil {
			fmt.Printf("stop read , err: %s\n", err)
		}
	}()
	go func() {
		err := this.asyncWrite(wg)
		if err != nil {
			fmt.Printf("stop async write , err: %s\n", err)
		}
	}()

	wg.Wait()
}
