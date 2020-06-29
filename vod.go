package do

func NewVod(engine *Engine) *Vod {
	return &Vod{engine: engine}
}

type Vod struct {
	process *Process
	engine  *Engine
}

func (p *Vod) fork(process *Process) *Vod {

	n := new(Vod)
	n.process = process
	n.engine = p.engine
	return n
}
