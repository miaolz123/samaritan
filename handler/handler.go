package handler

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/miaolz123/samaritan/config"
	"github.com/miaolz123/samaritan/constant"
)

type response struct {
	Success bool
	Message string
	Data    interface{}
}

type event struct{}

func (e event) OnSendHeader(ctx *rpc.HTTPContext) {
	ctx.Response.Header().Set("Access-Control-Allow-Headers", "Authorization")
}

// Server ...
func Server() {
	port := config.String("port")
	service := rpc.NewHTTPService()
	handler := struct {
		User      user
		Exchange  exchange
		Algorithm algorithm
		Trader    runner
		Log       logger
	}{}
	service.Event = event{}
	service.AddBeforeFilterHandler(func(request []byte, ctx rpc.Context, next rpc.NextFilterHandler) (response []byte, err error) {
		ctx.SetInt64("start", time.Now().UnixNano())
		httpContext := ctx.(*rpc.HTTPContext)
		if httpContext != nil {
			ctx.SetString("username", parseToken(httpContext.Request.Header.Get("Authorization")))
		}
		return next(request, ctx)
	})
	service.AddInvokeHandler(func(name string, args []reflect.Value, ctx rpc.Context, next rpc.NextInvokeHandler) (results []reflect.Value, err error) {
		name = strings.Replace(name, "_", ".", 1)
		results, err = next(name, args, ctx)
		spend := (time.Now().UnixNano() - ctx.GetInt64("start")) / 1000000
		spendInfo := ""
		if spend > 1000 {
			spendInfo = fmt.Sprintf("%vs", spend/1000)
		} else {
			spendInfo = fmt.Sprintf("%vms", spend)
		}
		log.Printf("%16s() spend %s", name, spendInfo)
		return
	})
	service.AddAllMethods(handler)
	http.Handle("/api", service)
	http.Handle("/", http.FileServer(http.Dir("web/dist")))
	log.Printf("Smaritan v%v running at http://localhost:%v\n", constant.Version, port)
	http.ListenAndServe(":"+port, nil)
}
