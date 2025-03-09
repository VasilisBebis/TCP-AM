// Package that implements the client's side logic of the program
package client

import (
	"fmt"
	_ "net"
)

// Header consists of the header fields used
// in the client's query message to the server
type Header struct {
	Op_code        byte   // operation code (0 -> A, 1 -> B, 2 -> C)
	Length         byte   // length of the data portion of the message in bytes (padding EXCLUDED)
	Transaction_id uint16 // unique identifier for the message (query & response)
}

// Message represents the full query message
type Message struct {
	Header Header
	Data   []byte // message data serialized as a byte array (padding INCLUDED)
}

func HelloClient() {
	fmt.Println("Hello From Client")
}
