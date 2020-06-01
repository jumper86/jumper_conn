package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/jumper86/jumper_conn"
	"github.com/jumper86/jumper_conn/def"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
	"github.com/jumper86/jumper_transform"
	jtd "github.com/jumper86/jumper_transform/def"
	jti "github.com/jumper86/jumper_transform/interf"
)

const addr = "localhost:8801"

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/ws_connect", ws_connect)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("listen and serve failed, err: %s\n", err)
		return
	}
}

//
func ws_connect(w http.ResponseWriter, r *http.Request) {

	wsOp := def.ConnOptions{
		MaxMsgSize:     def.MaxMsgSize,
		ReadTimeout:    def.ReadTimeout,
		WriteTimeout:   def.WriteTimeout,
		AsyncWriteSize: def.AsyncWriteSize,
		Side:           def.ServerSide,

		PingPeriod:       def.PingPeriod,
		PongWait:         def.PongWait,
		CloseGracePeriod: def.CloseGracePeriod,
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	//note: transform 可以只定义一个，他本身是线程安全对
	ts := jumper_transform.Newtransform()
	ts.AddOp(jtd.PacketBinary, nil)

	var h Handler

	jconn, err := jumper_conn.NewwsConn(wsConn, &wsOp, &h)
	if err != nil {
		return
	}

	h.Init(jconn, ts)

	fmt.Printf("local addr: %s, remote addr: %s\n", jconn.LocalAddr(), jconn.RemoteAddr())

	h.Set("radom_num", int(10000))
}

type Handler struct {
	interf.Conn
	jti.Transform
}

func (this *Handler) Init(conn interf.Conn, ts jti.Transform) {
	this.Conn = conn
	this.Transform = ts
	this.Run()
}

func (this *Handler) OnMessage(data []byte) error {
	defer util.TraceLog("handler.OnMessage")()
	fmt.Printf("handler get data: %v\n", data)

	//解出结构
	var msg jti.Message
	err := this.Execute(jtd.Backward, data, &msg)
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

	time.Sleep(1 * time.Second)
	this.Close()

	return nil
}

func (this *Handler) OnClose(err error) {
	defer util.TraceLog("handler.OnClose")()
}
