package do

import "strings"

type Headers struct {
	check map[string]bool
	list  []*Header
}

func (p *Headers) Add(headers ...*Header) {
	for _, v := range headers {
		name := strings.ToLower(v.name)
		if p.check == nil {
			p.check = make(map[string]bool)
		}
		if _, ok := p.check[name]; ok {
			continue
		}
		p.check[name] = true
		p.list = append(p.list, v)
	}
}

func (p *Headers) All() []*Header {
	return p.list
}
