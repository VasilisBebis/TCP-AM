// Package that implements the client's side logic of the program
package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	// "github.com/VasilisBebis/TCP-AM/pkg/server"
	"log"
	"math/rand/v2"
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

func (mess *Message) NewMessage(op byte, data any) {
	tr_id := rand.Uint32()
	bin_tr_id := new(bytes.Buffer)
	err := binary.Write(bin_tr_id, binary.BigEndian, tr_id)
	if err != nil {
		log.Println(err)
	}
	if bin_tr_id.Len() != 2 {
		log.Println(err)
	}
	bin_data := SerializeData(data)
	length := len(bin_data)
	_ = length
	// TODO: finish the message (add padding if necessary)
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
func (c *Client) ConnectTo(server_ip string, server_port string) {
	conn, err := net.Dial("tcp", server_ip+":"+server_port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to server ", server_ip)
	c.Conn = conn
}

// CloseConn closes the active connection of the client (if one exists)
func (c *Client) CloseConn() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}

// SerializeData creates a byte array of the given data
func SerializeData(data any) []byte {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
