package main

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/client"
	"github.com/VasilisBebis/TCP-AM/pkg/server"
)

func main() {
	mess := client.NewMessage(0, []int8{-3, 3, 1, 1, 1, 1})
	if mess == nil {
		return
	}

	server_ip := "127.0.0.1"
	server_port := "12345"

	c := client.NewClient()

	fmt.Println()
	c.ConnectTo(server_ip, server_port)
	c.SendMessage(*mess)
	buf := make([]byte, 256)
	fmt.Println(mess.Header.Transaction_id)

	for {
		_, err := c.Conn.Read(buf)
		if err == nil {
			message := server.DeserializeMessage(buf)
			result := server.DeserializeResult(message.Result, mess.Header.Op_code)
			fmt.Println(result)
			fmt.Println(message)

			break
		}
	}

	defer c.CloseConn()
}
