package interf

type Transform interface {
	AddOp(opType int, direct bool, params []interface{}) bool
	Execute(input interface{}, output interface{}) error
	Reset()
}
