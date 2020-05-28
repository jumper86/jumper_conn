package main

import (
	"errors"
	"fmt"
	"github.com/jumper86/jumper_conn/def"
	jc "github.com/jumper86/jumper_conn/impl/conn"
	jt "github.com/jumper86/jumper_conn/impl/transform/transform"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
	"net"
	"time"
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
			ts.AddOp(def.PacketBinary, nil)
			tcpOp := jc.NewtcpConnOptions(def.ServerSide, def.MaxMsgSize, def.ReadTimeout, def.WriteTimeout, def.AsyncWriteSize)
			jconn, err := jc.NewtcpConn(c, tcpOp, &h)
			if err != nil {
				fmt.Printf("new tcp conn failed. err: %s\n", err)
				return
			}
			h.Init(jconn, ts)

			fmt.Printf("local addr: %s, remote addr: %s\n", jconn.LocalAddr(), jconn.RemoteAddr())

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
			h.Set("radom_num", int(10000))

		}(conn)
	}
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
	var msg interf.Message
	err := this.Execute(def.Backward, data, &msg)
	if err != nil {
		fmt.Printf("transform failed, err: %s\n", err)
		return err
	}
	num := this.Get("radom_num").(int)
	fmt.Printf("num: %d\n", num)
	this.Del("radom_num")
	if n := this.Get("radom_num"); n != nil {
		fmt.Printf("delete failed")
		return errors.New("delete failed.")
	}

	fmt.Printf("type: %d, content: %s\n", msg.Type, msg.Content)

	time.Sleep(5 * time.Second)
	this.Close()
	return nil
}

func (this *Handler) OnClose(err error) {
	util.TraceLog("handler.OnClose")
}
