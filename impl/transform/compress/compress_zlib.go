package compress

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
	"io"
)

type compressOpZlib struct {
	direct bool
}

func NewcompressOpZlib(direct bool, params []interface{}) interf.CompressOp {
	var op compressOpZlib
	op.init(direct, params)
	return &op
}

func (self *compressOpZlib) init(direct bool, params []interface{}) bool {
	self.direct = direct
	return true
}

func (self *compressOpZlib) Operate(input interface{}, output interface{}) (bool, error) {

	if self.direct {
		tmpOutput, err := self.Compress(input.([]byte))
		if err != nil {
			fmt.Printf("pack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil

	} else {
		tmpOutput, err := self.Decompress(input.([]byte))
		if err != nil {
			fmt.Printf("unpack failed. err: %s", err)
			return false, err
		}
		*(output.(*[]byte)) = tmpOutput
		return true, nil
	}

	return true, nil
}

func (self *compressOpZlib) Compress(data []byte) ([]byte, error) {
	defer util.TraceLog("compressOpZlib.Compress")()
	var buf bytes.Buffer
	c := zlib.NewWriter(&buf)

	_, err := c.Write(data)
	if err != nil {
		fmt.Printf("compress err: %s", err)
		return nil, err
	}

	//!!!注意：若是上面使用　defer c.Close() 会导致在下面解压时出现错误：　decompress err: unexpected EOF
	//这里使用c.Close则不会
	//但是对于decompress 中的　reader 却可以defer close
	c.Close()
	return buf.Bytes(), nil
}

func (self *compressOpZlib) Decompress(data []byte) ([]byte, error) {

	defer util.TraceLog("compressOpZlib.Decompress")()
	nr := bytes.NewReader(data)
	dc, err := zlib.NewReader(nr)
	if err != nil {
		fmt.Printf("decompress err: %s", err)
		return nil, err
	}
	defer dc.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, dc)
	if err != nil {
		fmt.Printf("decompress err: %s", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
