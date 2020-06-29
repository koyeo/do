package do

func NewMySql(engine *Engine) *MySql {
	return &MySql{engine: engine}
}

type MySql struct {
	engine  *Engine
	process *Process
}

func (p *MySql) Exec(sql string, result interface{}) *Process {
	return new(Process)
}

func (p *MySql) fork(process *Process) *MySql {
	n := new(MySql)
	n.process = process
	n.engine = p.engine
	return n
}
