package do

import (
	"fmt"
	"strings"
)

func NewMySql(engine *Engine) *MySql {
	return &MySql{engine: engine}
}

type MySql struct {
	engine  *Engine
	process *Process
}

func (p *MySql) Select(database *MySqlDatabase, sql string, result interface{}) *Process {

	if p.process.isAbort {
		return p.process.pass()
	}

	t := fmt.Sprintf("%T", result)

	if !strings.HasPrefix(t, "*") {
		err := fmt.Errorf("result should be pointer, where is %s", t)
		return p.process.Abort(fmt.Sprintf("【Mysql】%s: ", database.name), err)
	}

	if t == "*int" || t == "*int32" || t == "*int64" {
		return p.count(database, sql, result)
	} else if strings.HasPrefix(t, "*[]") {
		return p.find(database, sql, result)
	}

	return p.get(database, sql, result)
}

func (p *MySql) count(database *MySqlDatabase, sql string, result interface{}) *Process {

	res, err := database.db.SQL(sql).Count()
	if err != nil {
		err = fmt.Errorf("exec count sql error: " + err.Error())
		return p.process.Abort(fmt.Sprintf("【Mysql】%s: ", database.name), err)
	}

	switch result.(type) {
	case *int:
		*result.(*int) = int(res)
	case *int32:
		*result.(*int32) = int32(res)
	case *int64:
		*result.(*int64) = res
	}

	return p.process
}

func (p *MySql) find(database *MySqlDatabase, sql string, result interface{}) *Process {

	err := database.db.SQL(sql).Find(result)
	if err != nil {
		err = fmt.Errorf("exec find sql error: " + err.Error())
		return p.process.Abort(fmt.Sprintf("【Mysql】%s", database.name), err)
	}

	return p.process
}

func (p *MySql) get(database *MySqlDatabase, sql string, result interface{}) *Process {

	_, err := database.db.SQL(sql).Get(result)
	if err != nil {
		err = fmt.Errorf("exec get sql error: " + err.Error())
		return p.process.Abort(fmt.Sprintf("【Mysql】%s", database.name), err)
	}

	return p.process
}

func (p *MySql) fork(process *Process) *MySql {
	n := new(MySql)
	n.process = process
	n.engine = p.engine
	return n
}
