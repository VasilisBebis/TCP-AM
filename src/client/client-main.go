package main

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/client"
)

func main() {
	mess := client.NewMessage(0, []int8{-5, 2, 3, 4})
	if mess == nil {
		return
	}

	server_ip := "127.0.0.1"
	server_port := "12345"

	c := client.NewClient()

	fmt.Println()
	c.ConnectTo(server_ip, server_port)
	c.SendMessage(*mess)
	defer c.CloseConn()
}
