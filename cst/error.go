package cst

import "errors"

const (
	ErrCodeConnClosed = iota
)

var (
	ErrConnClosed             = errors.New("conn is closed.")
	ErrConnUnexpectedClosed   = errors.New("conn is unexpected closed.")
	ErrCreateConnInvalidParam = errors.New("create conn invalid param.")
)
