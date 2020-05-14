package impl

import (
	"encoding/binary"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/jumper86/jumper_conn/interf"

	"github.com/jumper86/jumper_conn/cst"
)

type RawConn struct {
	conn        net.Conn
	closed      int32
	writeBuffer chan []byte
	closeChan   chan struct{}
	connOptions
	handler interf.Handler
}

func NewRawConn(conn net.Conn, co connOptions, handler interf.Handler) (interf.JConn, error) {
	err := checkOp(co, handler)
	if err != nil {
		return nil, err
	}
	rc := &RawConn{
		conn:        conn,
		closed:      0,
		writeBuffer: make(chan []byte, co.asyncWriteSize),
		closeChan:   make(chan struct{}),
		connOptions: co,
		handler:     handler,
	}
	return rc, nil
}

func (this *RawConn) Run() {

	go this.read()
	go this.asyncWrite()

}

func (this *RawConn) GetConn() net.Conn {
	return this.conn
}

func (this *RawConn) close(err error) {
	swapped := atomic.CompareAndSwapInt32(&this.closed, 0, 1)
	if !swapped {
		return
	}
	//todo: clean resource
	close(this.closeChan)
	close(this.writeBuffer)

	this.conn.Close()

	this.handler.OnClose(err)
}

func (this *RawConn) Close() {
	this.close(nil)
}

func (this *RawConn) IsClosed() bool {
	closed := atomic.LoadInt32(&this.closed)
	return closed == 1
}

func (this *RawConn) read() {

	var err error
readLoop:
	for {
		select {
		case <-this.closeChan:
			err = cst.ErrConnClosed
			break readLoop
		default:

			this.setReadDeadline(this.readTimeout)

			length := make([]byte, 4)
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

func (this *RawConn) setWriteDeadline(timeout int64) {
	if this.writeTimeout > 0 {
		this.conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
}

func (this *RawConn) setReadDeadline(timeout int64) {
	if this.readTimeout > 0 {
		this.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	}
}

func (this *RawConn) Write(data []byte) error {
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

func (this *RawConn) AsyncWrite(data []byte) (err error) {

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

func (this *RawConn) asyncWrite() error {

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
