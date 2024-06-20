package main

import (
	"log"

	"github.com/jrm0316/api-students/api"
)

func main() {
	server := api.NewServer()

	server.ConfigureRoutes()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
