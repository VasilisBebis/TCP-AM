package main

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/client"
)

func main() {
	server_ip := "127.0.0.1"
	server_port := "12345"

	c := client.NewClient()

	fmt.Println()
	c.ConnectTo(server_ip, server_port)
	defer c.CloseConn()
}
