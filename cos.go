package do

func NewCos(engine *Engine) *Cos {
	return &Cos{engine: engine}
}

type Cos struct {
	engine  *Engine
	process *Process
}

func (p *Cos) fork(process *Process) *Cos {
	n := new(Cos)
	n.process = process
	n.engine = p.engine
	return n
}
