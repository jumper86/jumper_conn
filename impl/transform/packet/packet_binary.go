package packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
)

type Message struct {
	Type    uint16
	Content []byte
}

type packetOpBinary struct {
	direct bool
}

func NewpacketOpBinary(direct bool, params []interface{}) interf.PacketOp {
	var op packetOpBinary
	op.init(direct, params)
	return &op
}

func (self *packetOpBinary) init(direct bool, params []interface{}) bool {
	self.direct = direct
	return true
}

func (self *packetOpBinary) Operate(input interface{}, output interface{}) (bool, error) {

	if self.direct {
		tmpOutput, err := self.Pack(input)
		if err != nil {
			fmt.Printf("pack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil

	} else {
		err := self.Unpack(input.([]byte), output)
		if err != nil {
			fmt.Printf("unpack failed. err: %s", err)
			return false, err
		}
		return true, nil
	}

	return true, nil
}

//此函数中需要检查入参是否为 string / []byte
func (*packetOpBinary) Pack(originData interface{}) ([]byte, error) {
	defer util.TraceLog("packetOpBinary.Pack")()
	msg, ok := originData.(*Message)
	if !ok {
		return nil, errors.New("invalid param type, use Message struct.")
	}

	rst := make([]byte, len(msg.Content)+2)
	binary.BigEndian.PutUint16(rst, msg.Type)
	copy(rst[2:], msg.Content)
	return rst, nil

}

func (*packetOpBinary) Unpack(packData []byte, obj interface{}) error {

	defer util.TraceLog("packetOpBinary.Unpack")()

	var msg *Message
	var ok bool
	if msg, ok = obj.(*Message); !ok {
		return errors.New("invalid param type, use Message struct.")
	}

	msg.Type = binary.BigEndian.Uint16(packData[:2])
	msg.Content = packData[2:]

	return nil
}
