package do

func NewHeader(name string, value string) *Header {
	return &Header{name: name, value: value}
}

type Header struct {
	name string
	value string
}

