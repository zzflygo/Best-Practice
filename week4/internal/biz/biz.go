package biz

import (
	"homework/week4/internal/data"
	"time"
)

type Greeter struct {
	Id      int64
	Message data.Dao
}

func NewGreeter(d data.Dao) Greeter {
	return Greeter{
		Id:      time.Now().Unix(),
		Message: d,
	}
}
