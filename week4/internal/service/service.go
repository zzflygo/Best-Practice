package service

import (
	"fmt"
	"homework/week4/internal/biz"
)

type Event struct {
	biz biz.Greeter
}

func NewEvent(b biz.Greeter) Event {
	return Event{
		biz: b,
	}
}

func (e *Event) Start() {
	fmt.Println("service start...")
}
