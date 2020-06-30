package do

func NewEngine(name string) *Engine {
	engine := &Engine{name: name}
	engine.initMySql()
	engine.initCos()
	engine.initRedis()
	engine.initAssert()
	engine.initVod()
	return engine
}

type Engine struct {
	Assert    *Assert
	Request   *Request
	MySql     *MySql
	Redis     *Redis
	Storage   *Storage
	Vod       *Vod
	Cos       *Cos
	name      string
	processes []*Process
}

func (p *Engine) InitRequest(host string, headers ...*Header) {
	p.Request = NewRequest(p)
	p.Request.host = host
	if len(headers) > 0 {
		p.Request.initHeaders()
		p.Request.headers.Add(headers...)
	}
}

func (p *Engine) initMySql() {
	p.MySql = NewMySql(p)
}

func (p *Engine) initRedis() {
	p.Redis = NewRedis(p)
}

func (p *Engine) initCos() {
	p.Cos = NewCos(p)
}

func (p *Engine) initVod() {
	p.Vod = NewVod(p)
}

func (p *Engine) initAssert() {
	p.Assert = NewAssert(p)
}

func (p *Engine) InitStorage(root string) {
	p.Storage = NewStorage(p)
}

func (p *Engine) Start(name string) *Process {
	r := newProcess(p, name)
	p.processes = append(p.processes, r)
	return r
}

func (p *Engine) Exec(routines int) {

}
