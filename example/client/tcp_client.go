package main

import (
	"encoding/binary"
	"fmt"
	"github.com/jumper86/jumper_conn/cst"
	jc "github.com/jumper86/jumper_conn/impl/conn"
	jt "github.com/jumper86/jumper_conn/impl/transform/transform"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
	"net"
	"sync"
)

const addr = "localhost:8801"

func main() {
	var conn net.Conn
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("dial failed, err: %s\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func(c net.Conn) {
		defer func() {
			fmt.Printf("============")
			wg.Done()
		}()

		var h Handler
		ts := jt.Newtransform()
		ts.AddOp(interf.PacketBinary, nil)
		tcpOp := jc.NewtcpConnOptions(cst.ClientSide, cst.MaxMsgSize, cst.ReadTimeout, cst.WriteTimeout, cst.AsyncWriteSize)
		jconn, err := jc.NewtcpConn(c, tcpOp, &h)
		if err != nil {
			fmt.Printf("new tcp conn failed. err: %s\n", err)
			return
		}
		h.Init(jconn, ts)

		fmt.Printf("local addr: %s, remote addr: %s\n", jconn.LocalAddr(), jconn.RemoteAddr())

		//send hello
		state := fmt.Sprintf("this is client %s, hello", jconn.LocalAddr())

		msg := &interf.Message{
			Type:    1,
			Content: []byte(state),
		}

		var output []byte
		err = h.Execute(interf.Forward, msg, &output)
		if err != nil {
			fmt.Printf("transform failed, err: %s\n", err)
			return
		}

		length := len(output)
		head := make([]byte, 4)
		binary.BigEndian.PutUint32(head, uint32(length))

		sendMsg := make([]byte, 0, cst.TcpHeadSize+length)
		sendMsg = append(sendMsg, head...)
		sendMsg = append(sendMsg, output...)

		h.Write(sendMsg)
		if err != nil {
			fmt.Printf("write failed, err: %s\n", err)
			return
		}

		//GetConn() net.Conn
		//Close()
		//IsClosed() bool
		//
		//Write(data []byte) error
		//AsyncWrite(data []byte) error
		//
		//
		//Set(string, interface{})
		//Get(string) interface{}
		//Del(string)

	}(conn)

	wg.Wait()
}

type Handler struct {
	interf.Conn
	interf.Transform
}

func (this *Handler) Init(conn interf.Conn, ts interf.Transform) {
	this.Conn = conn
	this.Transform = ts
	this.Run()
}

func (this *Handler) OnMessage(data []byte) error {
	util.TraceLog("handler.OnMessage")
	fmt.Printf("handler get data: %v\n", data)
	return nil
}

func (this *Handler) OnClose(err error) {
	util.TraceLog("handler.OnClose")
}