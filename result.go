package do

const (
	None    = 0
	Success = 1
	Failed  = 2
)

type Result struct {
	status int
	title  string
	params string
	result string
}
