package do

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"net/http"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

func NewRequest(engine *Engine) *Request {
	return &Request{
		engine:       engine,
		statusParser: defaultRequestStatusParser,
	}
}

type Request struct {
	process      *Process
	engine       *Engine
	host         string
	headers      *Headers
	statusParser func(res *http.Response) (data []byte, err error)
}

func (p *Request) initHeaders() {
	if p.headers == nil {
		p.headers = new(Headers)
	}
}

func (p *Request) fork(process *Process) *Request {
	n := new(Request)
	n.process = process
	n.engine = p.engine
	n.host = p.host
	n.headers = p.headers
	n.statusParser = p.statusParser
	return n
}

func (p *Request) url(path string) string {
	url := strings.TrimSuffix(p.host, "/") + path
	url = strings.ReplaceAll(url, "\u200b", "")
	return url
}

func (p *Request) Post(path string, data, result interface{}) *Process {

	url := p.url(path)

	if p.process.isAbort {
		return p.process.pass(fmt.Sprintf("【Request】POST %s", url))
	}

	err := p.send(POST, url, data, result)
	if err != nil {
		return p.process.abort(fmt.Sprintf("【Request】POST %s", url), "", err.Error())
	}
	return p.process
}

func (p *Request) Get(path string, result interface{}) *Process {
	if p.process.isAbort {
		return p.process.pass("")
	}

	return p.process
}

func (p *Request) Delete(path string, result interface{}) *Process {
	if p.process.isAbort {
		return p.process.pass("")
	}

	return p.process
}

func (p *Request) Put(path string, data, result interface{}) *Process {
	if p.process.isAbort {
		return p.process.pass("")
	}

	return p.process
}

func (p *Request) Header(name, value string) *Process {
	if p.process.isAbort {
		return p.process.pass("")
	}

	return p.process
}

func (p *Request) Exec(sql string, result interface{}) *Process {

	return p.process
}

func (p *Request) send(method, url string, data, result interface{}) (err error) {

	payload, err := p.parseData(data)
	if err != nil {
		return
	}

	client := &http.Client{
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return
	}
	if p.headers != nil {
		for _, v := range p.headers.All() {
			req.Header.Add(v.name, v.value)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		e := res.Body.Close()
		if err == nil {
			err = e
		}
	}()

	if p.statusParser == nil {
		err = fmt.Errorf("【Request】status parser is nil")
		return
	}

	output, err := p.statusParser(res)
	if err != nil {
		return
	}

	err = p.parseResult(output, result)
	if err != nil {
		return
	}

	return
}

func (p *Request) parseData(data interface{}) (res io.Reader, err error) {

	t := fmt.Sprintf("%T", data)
	if t != "string" {
		var r []byte
		r, err = json.Marshal(data)
		if err != nil {
			return
		}
		res = bytes.NewReader(r)
		return
	}

	res = strings.NewReader(data.(string))

	return
}

func (p *Request) parseResult(output []byte, result interface{}) (err error) {

	t := fmt.Sprintf("%T", result)

	if t == "*string" {
		r := result.(*string)
		*r = string(output)
		return
	}

	err = json.Unmarshal(output, &result)
	if err != nil {
		return
	}

	return
}

func defaultRequestStatusParser(res *http.Response) (data []byte, err error) {

	status, err := simplejson.NewFromReader(res.Body)
	if err != nil {
		return
	}

	code, err := status.Get("code").Int64()
	if err != nil {
		return
	}

	if code != 200 {
		err = fmt.Errorf("status code is not 200")
		return
	}

	data, err = status.Get("data").MarshalJSON()
	if err != nil {
		return
	}

	return
}
