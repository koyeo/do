package do

func NewRedis(engine *Engine) *Redis {
	return &Redis{engine: engine}
}

type Redis struct {
	process *Process
	engine  *Engine
}

func (p *Redis) Key(key string, value string) *Process {
	return p.process
}

func (p *Redis) fork(process *Process) *Redis {
	n := new(Redis)
	n.process = process
	n.engine = p.engine
	return n
}
