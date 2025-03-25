package main

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/client"
	"github.com/VasilisBebis/TCP-AM/pkg/messages"
)

func main() {
	mess := messages.NewQuery(0, []int8{-3, 3, 1, 1, 1, 1})
	if mess == nil {
		return
	}

	server_ip := "127.0.0.1"
	server_port := "12345"

	c := client.NewClient()

	fmt.Println()
	c.ConnectTo(server_ip, server_port)
	c.SendQuery(*mess)
	buf := make([]byte, 256)

	for {
		_, err := c.Conn.Read(buf)
		if err == nil {
			message := messages.DeserializeResponse(buf)
			result := messages.DeserializeResult(message.Result, mess.Header.Op_code)
			fmt.Println(result)
			fmt.Println(message)
			break
		}
	}

	defer c.CloseConn()
}
