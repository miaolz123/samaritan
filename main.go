package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/miaolz123/samaritan/handler"
)

type resp struct {
	Status  int64
	Message string
	Data    interface{}
}

func middlewareAuth(name string, args []reflect.Value, ctx rpc.Context, next rpc.NextInvokeHandler) (results []reflect.Value, err error) {
	fmt.Println("middlewareAuth")
	return next(name, args, ctx)
}

func hello(name string) (result reflect.Value, err error) {
	result = reflect.ValueOf(resp{
		Status: 200,
	})
	return
}

func main() {
	go handler.Run()
	service := rpc.NewHTTPService()
	service.AddBeforeFilterHandler(middlewareAuth)
	service.AddFunction("hello", hello)
	http.ListenAndServe(":9888", service)
}
