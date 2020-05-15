package interf

type Handler interface {
	Init(conn Conn, ts Transform) //初始化连接和转换
	OnMessage(data []byte) error  //在该函数实现中使用 transform 进行转换，而不是放在通信代码中
	OnClose(err error)
}
