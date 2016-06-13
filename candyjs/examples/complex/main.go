package main

import (
	"fmt"
	"os"

	"github.com/mcuadros/go-candyjs"
)

//go:generate candyjs import time
//go:generate candyjs import net/http
//go:generate candyjs import io/ioutil
//go:generate candyjs import github.com/gin-gonic/gin
func main() {
	script := os.Args[1]
	fmt.Printf("Executing %q\n", script)

	ctx := candyjs.NewContext()
	ctx.PevalFile(script)
}
