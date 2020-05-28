package def

import "github.com/jumper86/jumper_error"

const (
	ErrConnClosedCode           = 11011
	ErrConnUnexpectedClosedCode = 11012
	ErrInvalidConnParamCode     = 11013

	ErrGetExternalIpCode = 12011
	ErrGetMacAddrCode    = 12012
)

var (
	ErrConnClosed           = jumper_error.New(ErrConnClosedCode, "conn is closed.")
	ErrConnUnexpectedClosed = jumper_error.New(ErrConnUnexpectedClosedCode, "conn is unexpected closed.")
	ErrInvalidConnParam     = jumper_error.New(ErrInvalidConnParamCode, "create conn invalid param.")

	ErrGetExternalIp = jumper_error.New(ErrGetExternalIpCode, "get external ip failed.")
	ErrGetMacAddr    = jumper_error.New(ErrGetMacAddrCode, "get mac addr failed.")
)
