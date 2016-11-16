package main

import (
	"net/http"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/miaolz123/samaritan/handler"
)

func hello(name string) string {
	return "Hello " + name + "!"
}

func main() {
	go handler.Run()
	service := rpc.NewHTTPService()
	service.AddFunction("hello", hello)
	http.ListenAndServe(":9888", service)
}
