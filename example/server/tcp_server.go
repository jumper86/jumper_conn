package main

import (
	"fmt"
	"github.com/jumper86/jumper_conn/cst"
	jc "github.com/jumper86/jumper_conn/impl/conn"
	jt "github.com/jumper86/jumper_conn/impl/transform/transform"

	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
	"net"
)

const addr = "localhost:8801"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("listen failed, err: %s\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept failed, err: %s\n", err)
			return
		}

		go func(c net.Conn) {
			var h Handler
			ts := jt.Newtransform()
			ts.AddOp(interf.PacketBinary, nil)
			tcpOp := jc.NewtcpConnOptions(cst.ServerSide, cst.MaxMsgSize, cst.ReadTimeout, cst.WriteTimeout, cst.AsyncWriteSize)
			tcpConn, err := jc.NewtcpConn(c, tcpOp, &h)
			if err != nil {
				fmt.Printf("new tcp conn failed. err: %s\n", err)
				return
			}
			h.Init(tcpConn, ts)

			fmt.Printf("local addr: %s, remote addr: %s\n", tcpConn.LocalAddr(), tcpConn.RemoteAddr())

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
	}
}

type Handler struct {
	conn interf.Conn
	ts   interf.Transform
}

func (this *Handler) Init(conn interf.Conn, ts interf.Transform) {
	this.conn = conn
	this.ts = ts
	this.conn.Run()
}

func (this *Handler) OnMessage(data []byte) error {
	util.TraceLog("handler.OnMessage")
	fmt.Printf("handler get data: %v\n", data)
	return nil
}

func (this *Handler) OnClose(err error) {
	util.TraceLog("handler.OnClose")
}
