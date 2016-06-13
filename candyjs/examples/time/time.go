package main

import (
	"time"

	"github.com/mcuadros/go-candyjs"
)

func main() {
	ctx := candyjs.NewContext()
	ctx.PushGlobalGoFunction("date", time.Date)
	ctx.PushGlobalGoFunction("now", time.Now)
	ctx.PushGlobalProxy("UTC", time.UTC)

	ctx.EvalString(`
        future = date(2015, 10, 21, 4, 29 ,0, 0, UTC)

        print("Back to the Future day is on: " + future.sub(now()) + " nsecs!")
    `)
}
