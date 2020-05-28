package encrypt

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
)

type encryptOpSha1 struct {
}

func NewencryptOpSha1(params []interface{}) interf.EncryptOp {
	var op encryptOpSha1
	op.init(params)
	return &op
}

func (self *encryptOpSha1) init(params []interface{}) bool {
	return true
}

func (self *encryptOpSha1) Operate(direct int8, input interface{}, output interface{}) (bool, error) {

	if direct == def.Forward {
		tmpOutput, err := self.Encrypt(input.([]byte))
		if err != nil {
			fmt.Printf("pack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil

	} else {
		fmt.Println("sha1 couldn't decrypt.")
		return false, errors.New("sha1 couldn't decrypt.")
	}

	return true, nil
}

func (*encryptOpSha1) Encrypt(data []byte) ([]byte, error) {
	defer util.TraceLog("encryptOpSha1.Encrypt")()
	r := sha1.Sum(data)
	rst := r[:]
	return rst, nil

}

func (*encryptOpSha1) Decrypt(data []byte) ([]byte, error) {
	defer util.TraceLog("encryptOpSha1.Decrypt")()
	return nil, errors.New("sha1 couldn't decrypt.")
}
