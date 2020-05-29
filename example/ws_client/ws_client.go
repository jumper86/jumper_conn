package main

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/jumper86/jumper_conn"
	"github.com/jumper86/jumper_conn/def"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
)

const addr = "localhost:8801"

func main() {

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws_connect"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Printf("dial failed, err: %s\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func(c *websocket.Conn) {
		defer func() {
			fmt.Printf("============")
			wg.Done()
		}()

		var h Handler
		ts := jumper_conn.Newtransform()
		ts.AddOp(def.PacketJson, nil)

		wsOp := def.ConnOptions{
			MaxMsgSize:     def.MaxMsgSize,
			ReadTimeout:    def.ReadTimeout,
			WriteTimeout:   def.WriteTimeout,
			AsyncWriteSize: def.AsyncWriteSize,
			Side:           def.ClientSide,

			PingPeriod:       def.PingPeriod,
			PongWait:         def.PongWait,
			CloseGracePeriod: def.CloseGracePeriod,
		}

		jconn, err := jumper_conn.NewwsConn(c, &wsOp, &h)
		if err != nil {
			fmt.Printf("new tcp conn failed. err: %s\n", err)
			return
		}
		h.Init(jconn, ts)

		fmt.Printf("local addr: %s, remote addr: %s\n", jconn.LocalAddr(), jconn.RemoteAddr())

		//send hello
		str := fmt.Sprintf("this is tcp_client %s, hello", jconn.LocalAddr())
		//
		type State struct {
			Type    int64  `json:"type"`
			Content string `json:"content"`
		}

		msg := State{Type: 1, Content: str}

		var output []byte
		err = h.Execute(def.Forward, msg, &output)
		if err != nil {
			fmt.Printf("transform failed, err: %s\n", err)
			return
		}

		sendMsg := output
		fmt.Printf("sendMsg: %v\n", sendMsg)

		//length := len(output)
		//head := make([]byte, 4)
		//binary.BigEndian.PutUint32(head, uint32(length))
		//
		//sendMsg := make([]byte, 0, def.TcpHeadSize+length)
		//sendMsg = append(sendMsg, head...)
		//sendMsg = append(sendMsg, output...)

		h.Write(sendMsg)
		if err != nil {
			fmt.Printf("write failed, err: %s\n", err)
			return
		}

	}(c)

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
