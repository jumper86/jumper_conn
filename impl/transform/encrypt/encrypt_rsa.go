package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
)

type encryptOpRsa struct {
	rsaPublicKeyRemote []byte // 来自对端生成的公钥，用于加密
	rsaPrivateKeyLocal []byte // 来自本端生成的私钥
	direct             bool
}

func NewencryptOpRsa(direct bool, params []interface{}) interf.EncryptOp {
	var op encryptOpRsa
	op.init(direct, params)
	return &op
}

func (self *encryptOpRsa) init(direct bool, params []interface{}) bool {
	if len(params) != 2 {
		fmt.Printf("invalid param count.")
		return false
	}

	var ok bool
	self.rsaPublicKeyRemote, ok = params[0].([]byte)
	if !ok {
		fmt.Printf("invalid param type ")
		return false
	}
	self.rsaPrivateKeyLocal, ok = params[1].([]byte)
	if !ok {
		fmt.Printf("invalid param type ")
		return false
	}

	if len(self.rsaPrivateKeyLocal) > 2048 {
		fmt.Printf("invalid param type ")
		return false
	}

	self.direct = direct

	return true
}

func (self *encryptOpRsa) Operate(input interface{}, output interface{}) (bool, error) {

	if self.direct {
		tmpOutput, err := self.Encrypt(input.([]byte))
		if err != nil {
			fmt.Printf("pack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil

	} else {
		tmpOutput, err := self.Decrypt(input.([]byte))
		if err != nil {
			fmt.Printf("unpack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil
	}

	return true, nil
}

func (self *encryptOpRsa) Encrypt(data []byte) ([]byte, error) {

	defer util.TraceLog("encryptOpRsa.Encrypt")()
	block, _ := pem.Decode(self.rsaPublicKeyRemote)
	if block == nil {
		return nil, errors.New("public key error")

	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err

	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

func (self *encryptOpRsa) Decrypt(data []byte) ([]byte, error) {

	defer util.TraceLog("encryptOpRsa.Decrypt")()
	block, _ := pem.Decode(self.rsaPrivateKeyLocal)
	if block == nil {
		return nil, errors.New("private key error!")

	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err

	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}
