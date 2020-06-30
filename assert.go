package do

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"strings"
)

func NewAssert(engine *Engine) *Assert {
	return &Assert{engine: engine}
}

type Assert struct {
	process *Process
	engine  *Engine
}

func (p *Assert) fork(process *Process) *Assert {
	n := new(Assert)
	n.process = process
	n.engine = p.engine
	return n
}

func (p *Assert) Equal(x, y interface{}, message ...string) *Process {
	if p.process.isAbort {
		return p.process.pass()
	}

	msg := strings.Join(message, "")
	if !cmp.Equal(x, y) {
		err := fmt.Errorf("%+v is not equal %+v", x, y)
		return p.process.Abort(msg, err)
	}

	return p.process.next("【Assert】", "", "")
}

func (p *Assert) NotEqual(x, y interface{}, message ...string) *Process {
	if p.process.isAbort {
		return p.process.pass()
	}

	msg := strings.Join(message, "")
	if cmp.Equal(x, y) {
		err := fmt.Errorf("%+v is equal %+v", x, y)
		return p.process.Abort(msg, err)
	}

	return p.process.next("【Assert】", "", "")
}
