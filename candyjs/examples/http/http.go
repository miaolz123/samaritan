package main

import (
	"io"
	"net/http"

	"github.com/mcuadros/go-candyjs"
)

func main() {
	ctx := candyjs.NewContext()
	ctx.PushGlobalGoFunction("handleFunc", http.HandleFunc)
	ctx.PushGlobalGoFunction("listenAndServe", http.ListenAndServe)
	ctx.PushGlobalGoFunction("writeString", io.WriteString)

	ctx.EvalString(`
        handler = function(writer, request) {
            writeString(writer, "Hello from CandyJS!")
        }

        handleFunc("/", CandyJS.proxy(handler))
        listenAndServe(":8000", null)
    `)
}
