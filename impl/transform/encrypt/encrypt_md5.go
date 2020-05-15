package encrypt

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
)

type encryptOpMd5 struct {
	direct bool
}

func NewencryptOpMd5(direct bool, params []interface{}) interf.EncryptOp {
	var op encryptOpMd5
	op.init(direct, params)
	return &op
}

func (self *encryptOpMd5) init(direct bool, params []interface{}) bool {
	self.direct = direct
	return true
}

func (self *encryptOpMd5) Operate(input interface{}, output interface{}) (bool, error) {

	if self.direct {
		tmpOutput, err := self.Encrypt(input.([]byte))
		if err != nil {
			fmt.Printf("pack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil

	} else {
		fmt.Println("md5 couldn't decrypt.")
		return false, errors.New("md5 couldn't decrypt.")
	}

	return true, nil
}

func (*encryptOpMd5) Encrypt(data []byte) ([]byte, error) {
	defer util.TraceLog("encryptOpMd5.Encrypt")()
	r := md5.Sum(data)
	rst := r[:]
	return rst, nil

}

func (*encryptOpMd5) Decrypt(data []byte) ([]byte, error) {
	defer util.TraceLog("encryptOpMd5.Decrypt")()
	return nil, errors.New("md5 couldn't decrypt.")
}
