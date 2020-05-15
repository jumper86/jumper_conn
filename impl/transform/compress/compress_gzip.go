package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/jumper86/jumper_conn/interf"
	"github.com/jumper86/jumper_conn/util"
	"io"
)

type compressOpGzip struct {
}

func NewcompressOpGzip(params []interface{}) interf.CompressOp {
	var op compressOpGzip
	op.init(params)
	return &op
}

func (self *compressOpGzip) init(params []interface{}) bool {
	return true
}

func (self *compressOpGzip) Operate(direct int8, input interface{}, output interface{}) (bool, error) {

	if direct == interf.Forward {
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

func (self *compressOpGzip) Compress(data []byte) ([]byte, error) {

	defer util.TraceLog("compressOpGzip.Compress")()
	var buf bytes.Buffer
	c := gzip.NewWriter(&buf)

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

func (self *compressOpGzip) Decompress(data []byte) ([]byte, error) {

	defer util.TraceLog("compressOpGzip.Decompress")()
	nr := bytes.NewReader(data)
	dc, err := gzip.NewReader(nr)
	if err != nil {
		fmt.Printf("decompress 1 err: %s", err)
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
