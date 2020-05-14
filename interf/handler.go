package interf

type Handler interface {
	OnMessage(data []byte) error //在该函数实现中使用 transform 进行转换，而不是放在通信代码中
	OnClose(err error)
}
