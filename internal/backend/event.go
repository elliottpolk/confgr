package backend

import (
	"encoding/json"
	"time"
)

type Op int

// ops
const (
	SetOp Op = iota
	GetOp
	RemoveOp
	ListKeysOp
)

type Event map[string]interface{}

func NewEvent(o Op, r, d, u string) Event {
	return Event{
		"op": o,                     // operation performed
		"r":  r,                     // repo the op was performed on
		"d":  d,                     // any data affected in the op
		"u":  u,                     // user that performed the op
		"t":  time.Now().UnixNano(), // timestamp the even occured
	}
}

func (e Event) Fail() Event {
	_e := Event{"failed": true}
	for k, v := range e {
		_e[k] = v
	}
	return _e
}

func (e Event) String() string {
	out, _ := json.Marshal(e)
	return string(out)
}
