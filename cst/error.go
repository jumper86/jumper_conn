package cst

import "errors"

var (
	ErrConnClosed             = errors.New("conn is closed.")
	ErrConnUnexpectedClosed   = errors.New("conn is unexpected closed.")
	ErrCreateConnInvalidParam = errors.New("create conn invalid param.")
)

var (
	ErrGetExternalIp = errors.New("get external ip failed.")
	ErrGetMacAddr    = errors.New("get mac addr failed.")
)
