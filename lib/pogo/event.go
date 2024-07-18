package pogo

import (
	"sync"

	"github.com/goccy/go-json"
)

var (
	eventPool = sync.Pool{
		New: func() any {
			return &Event{}
		},
	}

	AcqEvent = func() *Event {
		event := eventPool.Get().(*Event)
		return event
	}

	AcqEventWithId = func(id int64) *Event {
		event := AcqEvent()
		event.Id = id
		return event
	}
)

type Event struct {
	Id   int64 `json:"id"`
	Data Param `json:"data"`
}

func (e *Event) WithData(data Param) *Event {
	e.Data = data
	return e
}

func (e *Event) Marshal() []byte {
	data, _ := json.Marshal(e)
	return data
}

func (e *Event) Close() {
	e.reset()
	eventPool.Put(e)
}

func (e *Event) reset() {
	e.Id = 0
	e.Data = nil
}
