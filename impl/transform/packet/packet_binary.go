package packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jumper86/jumper_conn/def"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
)

type packetOpBinary struct {
}

func NewpacketOpBinary(params []interface{}) interf.PacketOp {
	var op packetOpBinary
	op.init(params)
	return &op
}

func (self *packetOpBinary) init(params []interface{}) bool {
	return true
}

func (self *packetOpBinary) Operate(direct int8, input interface{}, output interface{}) (bool, error) {

	if direct == def.Forward {
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
	msg, ok := originData.(*interf.Message)
	if !ok {
		return nil, errors.New("invalid param type, use interf.Message struct.")
	}

	rst := make([]byte, len(msg.Content)+2)
	binary.BigEndian.PutUint16(rst, msg.Type)
	copy(rst[2:], msg.Content)
	return rst, nil

}

func (*packetOpBinary) Unpack(packData []byte, obj interface{}) error {

	defer util.TraceLog("packetOpBinary.Unpack")()

	var msg *interf.Message
	var ok bool
	if msg, ok = obj.(*interf.Message); !ok {
		return errors.New("invalid param type, use interf.Message struct.")
	}

	msg.Type = binary.BigEndian.Uint16(packData[:2])
	msg.Content = packData[2:]

	return nil
}
