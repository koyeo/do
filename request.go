package do

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
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
	host := strings.ReplaceAll(p.host, "\u200b", "")
	path = strings.ReplaceAll(path, "\u200b", "")
	url := strings.TrimSuffix(host, "/") + "/" + strings.TrimPrefix(path, "/")
	return url
}

func (p *Request) Post(path string, data, result interface{}) *Process {

	url := p.url(path)

	if p.process.isAbort {
		return p.process.pass()
	}

	err := p.send(POST, &url, data, result)
	if err != nil {
		return p.process.Abort(fmt.Sprintf("【Request】POST %s", url), err)
	}
	return p.process
}

func (p *Request) Get(path string, result interface{}) *Process {

	url := p.url(path)

	if p.process.isAbort {
		return p.process.pass()
	}

	err := p.send(GET, &url, "", result)
	if err != nil {
		return p.process.Abort(fmt.Sprintf("【Request】GET %s", url), err)
	}
	return p.process
}

func (p *Request) Delete(path string, result interface{}) *Process {

	url := p.url(path)

	if p.process.isAbort {
		return p.process.pass()
	}

	err := p.send(DELETE, &url, "", result)
	if err != nil {
		return p.process.Abort(fmt.Sprintf("【Request】POST %s", url), err)
	}
	return p.process
}

func (p *Request) Put(path string, data, result interface{}) *Process {

	url := p.url(path)

	if p.process.isAbort {
		return p.process.pass()
	}

	err := p.send(PUT, &url, data, result)
	if err != nil {
		return p.process.Abort(fmt.Sprintf("【Request】POST %s", url), err)
	}
	return p.process
}

func (p *Request) Route(route Route, data, result interface{}) *Process {

	url := p.url(route.Path())
	method := strings.ToUpper(route.Method())

	if p.process.isAbort {
		return p.process.pass()
	}

	err := p.send(method, &url, data, result)

	if err != nil {
		return p.process.Abort(fmt.Sprintf("【Request】%s %s", method, url), err)
	}

	return p.process
}

func (p *Request) Header(name, value string) *Process {
	if p.process.isAbort {
		return p.process.pass()
	}

	p.initHeaders()
	p.headers.Add(NewHeader(name, value))

	return p.process
}

func (p *Request) PrintHeaders() {
	if p.headers != nil {
		for _, v := range p.headers.All() {
			fmt.Printf("%s: %s\n", v.name, v.value)
		}
	} else {
		fmt.Println("request headers is empty")
	}
}

func (p *Request) send(method string, url *string, data, result interface{}) (err error) {

	var payload io.Reader

	if method == POST || method == PUT {
		payload, err = p.parseData(data)
		if err != nil {
			err = fmt.Errorf("parse data error: " + err.Error())
			return
		}
	}

	client := &http.Client{
	}
	req, err := http.NewRequest(method, *url, payload)

	if method == GET || method == DELETE {
		err = setRequestParams(req, data)
		if err != nil {
			return
		}
	}

	if req != nil {
		*url = req.URL.String()
	}

	if err != nil {
		err = fmt.Errorf("new request error: " + err.Error())
		return
	}
	if p.headers != nil {
		for _, v := range p.headers.All() {
			req.Header.Add(v.name, v.value)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("send request error: " + err.Error())
		return
	}
	defer func() {
		e := res.Body.Close()
		if err == nil && e != nil {
			err = fmt.Errorf("close request error: " + e.Error())
		}
	}()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("http code is not 200")
		return
	}

	if p.statusParser == nil {
		err = fmt.Errorf("status parser is nil")
		return
	}

	output, err := p.statusParser(res)
	if err != nil {
		err = fmt.Errorf("parse status error: " + err.Error())
		return
	}

	if result != nil {
		err = p.parseResult(output, result)
		if err != nil {
			err = fmt.Errorf("parse result error: " + err.Error())
			return
		}
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("read body error: " + err.Error())
		return
	}
	defer func() {
		if err != nil {
			content := string(body)
			if content != "" {
				err = fmt.Errorf(content)
			} else {
				err = fmt.Errorf("response body is empty")
			}

		}
	}()

	status, err := simplejson.NewJson(body)
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

func setRequestParams(req *http.Request, data interface{}) (err error) {

	defer func() {
		e := recover()
		if err != nil {
			err = fmt.Errorf("ser request params error: %+v", e)
		}
	}()

	j, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("json marshal error: %+v", err)
		return
	}

	_params, err := simplejson.NewJson(j)
	if err != nil {
		err = fmt.Errorf("new simple json error: %+v", err)
		return
	}

	params, err := _params.Map()
	if err != nil {
		err = fmt.Errorf("get params map error: %+v", err)
		return
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, fmt.Sprintf("%+v", v))
	}
	req.URL.RawQuery = q.Encode()

	return
}
