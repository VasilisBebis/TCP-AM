package main

import (
	"encoding/hex"
	"fmt"

	"github.com/VasilisBebis/TCP-AM/pkg/client"
)

func main() {
	mess := client.NewMessage(0, []int8{1, 2, 3, 4})
	if mess == nil {
		fmt.Println("Test")
		return
	}
	fmt.Printf("Message: %#v\n", mess)
	fmt.Printf("%v\n", mess.SerializeMessage())

	str := hex.EncodeToString(mess.SerializeMessage())
	fmt.Println(str)
	// server_ip := "127.0.0.1"
	// server_port := "12345"
	//
	// c := client.NewClient()
	//
	// fmt.Println()
	// c.ConnectTo(server_ip, server_port)
	// defer c.CloseConn()
}
