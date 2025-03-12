// Package that implements the client's side logic of the program
package client

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/server"
	"log"
	"net"
)

// Header consists of the header fields used
// in the client's query message to the server
type Header struct {
	Op_code        byte    // operation code (0 -> A, 1 -> B, 2 -> C)
	Length         byte    // length of the data portion of the message in bytes (padding EXCLUDED)
	Transaction_id [2]byte // unique identifier for the message (query & response)
}

// Message represents the full query message
type Message struct {
	Header  Header
	Data    []byte // message data serialized as a byte array
	Padding []byte // used to make the message 32-bit aligned (empty if not needed)
}

type Client struct {
	Conn net.Conn // client's connector
}

// NewClient creates an object of type Client and return it
func NewClient() *Client {
	client := Client{Conn: nil}
	return &client

}

// ConnectTo creates a connection between the client and the given server
func (c *Client) ConnectTo(s server.Server) {

	connc, err := net.Dial("tcp", "localhost:"+s.Port)
	if err != nil {
		log.Fatal(err)
	}
	c.Conn = connc
}

// CloseConn closes the active connection of the client (if one exists)
func (c *Client) CloseConn() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}

func HelloClient() {
	fmt.Println("Hello From Client")
}
