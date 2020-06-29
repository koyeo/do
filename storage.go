package do

func NewStorage(engine *Engine) *Storage {
	return &Storage{engine: engine}
}

type Storage struct {
	engine  *Engine
	process *Process
	root    string
}

func (p *Storage) fork(process *Process) *Storage {
	n := new(Storage)
	n.process = process
	n.engine = p.engine
	n.root = p.root
	return n
}

