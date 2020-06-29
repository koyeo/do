package do

import (
	"github.com/ttacon/chalk"
	"log"
)

func newProcess(engine *Engine, name string) *Process {
	r := &Process{
		name: name,
	}
	if engine.Request != nil {
		r.Request = engine.Request.fork(r)
	}
	if engine.MySql != nil {
		r.MySql = engine.MySql.fork(r)
	}
	if engine.Redis != nil {
		r.Redis = engine.Redis.fork(r)
	}
	if engine.Storage != nil {
		r.Storage = engine.Storage.fork(r)
	}
	if engine.Vod != nil {
		r.Vod = engine.Vod.fork(r)
	}
	if engine.Cos != nil {
		r.Cos = engine.Cos.fork(r)
	}

	return r
}

type Process struct {
	Request  *Request
	MySql    *MySql
	Redis    *Redis
	Storage  *Storage
	Vod      *Vod
	Cos      *Cos
	name     string
	async    bool
	isAbort  bool
	requests []*Request
	results  []*Result
}

func (p *Process) addResult(result *Result) {
	p.results = append(p.results, result)
}

func (p *Process) abort(title, params, result string) *Process {
	p.isAbort = true
	p.addResult(&Result{
		status: Failed,
		title:  title,
		params: params,
		result: result,
	})
	log.Println(chalk.Red.Color(chalk.Bold.TextStyle(title)), chalk.Red.Color(result))
	//debug.PrintStack()
	return p
}

func (p *Process) next(title, params, result string) *Process {

	p.addResult(&Result{
		status: Success,
		title:  title,
		params: params,
		result: result,
	})

	return p
}

func (p *Process) pass(title string) *Process {

	p.addResult(&Result{
		status: None,
		title:  title,
	})

	return p
}
