package do

type Route interface {
	Method() string
	Path() string
}
