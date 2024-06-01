package main

import (
	"log"
	"otel/internal/infra/webserver"
)

func main() {
	if err := webserver.StartServer(webserver.NewServerB()); err != nil {
		log.Panic(err)
	}
}
