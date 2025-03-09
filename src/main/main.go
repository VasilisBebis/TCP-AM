package main

import (
	"github.com/VasilisBebis/TCP-AM/src/client"
	"github.com/VasilisBebis/TCP-AM/src/server"
	"os"
	// "errors"
)

func main() {
	if len(os.Args) < 2 {
		panic("No arguments provided!")
	}
	program := os.Args[1]

	if program == "server" {
		// _, err := os.Stat("../client/client.go")
		// if errors.Is(err, os.ErrNotExist) {
		//
		// }
		serv := server.NewServer()
		serv.OpenServer()

		server.HelloServer()
	} else if program == "client" {
		client.HelloClient()
	} else {
		panic("Given program not declared")
	}

}
