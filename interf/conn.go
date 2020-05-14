package interf

import (
	"net"
)

type JConn interface {
	GetConn() net.Conn
	Close()
	IsClosed() bool

	Run()

	Write(data []byte) error
	AsyncWrite(data []byte) error
}

type Handler interface {
	OnMessage(data []byte) error
	OnClose(err error)
}

//
//type Packet struct{
//
//}
//
//type Protocol interface {
//	Marshale()
//	UnMarshale(reader io.Reader, buffer []byte) error
//}
