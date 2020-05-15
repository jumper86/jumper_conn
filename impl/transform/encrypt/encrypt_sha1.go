package encrypt

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
)

type encryptOpSha1 struct {
	direct bool
}

func NewencryptOpSha1(direct bool, params []interface{}) interf.EncryptOp {
	var op encryptOpSha1
	op.init(direct, params)
	return &op
}

func (self *encryptOpSha1) init(direct bool, params []interface{}) bool {
	self.direct = direct
	return true
}

func (self *encryptOpSha1) Operate(input interface{}, output interface{}) (bool, error) {

	if self.direct {
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
