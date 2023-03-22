package wfile

type EventCode int

const (
	CHANGE EventCode = iota
	NOCHANGE
	ERROR
)

// Event represents the event data when changes occur.
type Event struct {
	name  string
	code  EventCode
	path  string
	error error
}
