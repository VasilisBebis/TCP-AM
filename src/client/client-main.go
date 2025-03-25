package main

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/client"
	"github.com/VasilisBebis/TCP-AM/pkg/messages"
)

func main() {
	set1 := []uint16{2, 4, 6}
	set2 := []uint16{1, 3, 4}
	mess := messages.NewQuery(2, append(set2, set1...))
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
			if message.Header.Response_code != messages.Ok {
				fmt.Println(messages.Response_message[message.Header.Response_code])
			} else {
				result := messages.DeserializeResult(message.Result, mess.Header.Op_code)
				fmt.Println("The result of the operation is:", result)
			}
			break
		}
	}

	defer c.CloseConn()
}
