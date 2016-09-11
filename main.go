package main

import "github.com/miaolz123/samaritan/handler"

func main() {
	handler.Server.Listen(":9806")
}
