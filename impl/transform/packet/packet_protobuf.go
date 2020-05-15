package packet

import (
	"errors"
	"github.com/jumper86/jumper_conn/interf"

	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jumper86/jumper_conn/util"
)

type packetOpProtobuf struct {
}

func NewpacketOpProtobuf(params []interface{}) interf.PacketOp {
	var op packetOpProtobuf
	op.init(params)
	return &op
}

func (self *packetOpProtobuf) init(params []interface{}) bool {
	return true
}

func (self *packetOpProtobuf) Operate(direct int8, input interface{}, output interface{}) (bool, error) {

	if direct == interf.Forward {
		tmpOutput, err := self.Pack(input)
		if err != nil {
			fmt.Printf("pack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil

	} else {
		//fmt.Printf("output: %v", output)
		err := self.Unpack(input.([]byte), output)
		if err != nil {
			fmt.Printf("unpack failed. err: %s", err)
			return false, err
		}
		return true, nil
	}

	return true, nil
}

func (*packetOpProtobuf) Pack(originData interface{}) ([]byte, error) {
	//此处需要将interface{} -> proto.Message， 使用类型断言即可
	defer util.TraceLog("packetOpProtobuf.Pack")()
	data, ok := originData.(proto.Message)
	if !ok {
		return nil, errors.New("param not implement interface proto.Message.")
	}

	return proto.Marshal(data)
}

func (*packetOpProtobuf) Unpack(packData []byte, obj interface{}) error {

	defer util.TraceLog("packetOpProtobuf.Unpack")()
	decodedData, ok := obj.(proto.Message)
	if !ok {
		return errors.New("param not implement interface proto.Message.")
	}
	err := proto.Unmarshal(packData, decodedData)
	return err

}
