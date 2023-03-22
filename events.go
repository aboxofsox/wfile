package wfile

type EventCode int

const (
	CHANGE EventCode = iota
	NOCHANGE
	ERROR
)

type Event struct {
	name  string
	code  EventCode
	path  string
	error error
}
