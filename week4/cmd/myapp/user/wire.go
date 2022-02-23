//go:build wireinject
// +build wireinject

package main //需要和上面的自定义标志,空一格不然会失败

import (
	"github.com/google/wire"
	"homework/week4/internal/biz"
	"homework/week4/internal/data"
	"homework/week4/internal/service"
)

func InitService(sql string) service.Event {
	// wire.Build 传参只要函数名 --不能data.NewDao()带括号
	wire.Build(data.NewDao, biz.NewGreeter, service.NewEvent)
	// return需要返回一个结构体实例
	return service.Event{}
}
