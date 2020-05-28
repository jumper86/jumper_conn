package jumper_conn

import (
	"github.com/jumper86/jumper_conn/impl/transform/transform"
	"github.com/jumper86/jumper_conn/interf"
)

func Newtransform() interf.Transform {
	var tf transform.Transform
	return &tf
}
