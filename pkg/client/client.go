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
	Op_code        byte   // operation code (0 -> A, 1 -> B, 2 -> C)
	Length         byte   // length of the data portion of the message in bytes (padding EXCLUDED)
	Transaction_id []byte // unique identifier for the message (query & response)
}

// Message represents the full query message
type Message struct {
	Header  Header
	Data    []byte // message data serialized as a byte array
	Padding []byte // used to make the message 32-bit aligned (empty if not needed)
}

// Generates a new Message object that includes the given data
func NewMessage(op byte, data any) *Message {
	tr_id := uint16(rand.Uint32() >> 16)
	bin_tr_id := new(bytes.Buffer)
	err := binary.Write(bin_tr_id, binary.BigEndian, tr_id)

	if err != nil {
		log.Println(err)
	}
	if bin_tr_id.Len() != 2 {
		log.Println(err)
		return nil
	}

	bin_data, err := SerializeData(data)
	if err != nil {
		return nil
	}
	length := len(bin_data)
	padding_size := (4 - length%4) % 4
	padding := make([]byte, padding_size)
	header := Header{Op_code: op, Length: byte(length), Transaction_id: bin_tr_id.Bytes()}
	message := Message{Header: header, Data: bin_data, Padding: padding}
	return &message
}

// SerializeMessage packs an object of type [Message] to a byte stream (big endian)
func (m *Message) SerializeMessage() []byte {
	var message_bytes []byte

	message_bytes = append(message_bytes, []byte{m.Header.Op_code, m.Header.Length}...)
	message_bytes = append(message_bytes, m.Header.Transaction_id[:]...)
	message_bytes = append(message_bytes, m.Data...)
	if len(m.Padding) != 0 {
		message_bytes = append(message_bytes, m.Padding...)
	}
	return message_bytes
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
		log.Println(err)
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
func SerializeData(data any) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch v := data.(type) {
	case []int8:
		for _, i := range v {
			binary.Write(buf, binary.BigEndian, i)
		}
	case []uint8:
		return []byte(v), nil
	case []uint16:
		for _, i := range v {
			binary.Write(buf, binary.BigEndian, i)
		}
	default:
		return nil, fmt.Errorf("Type %T is not supported!", v)

	}
	// err := binary.Write(buf, binary.BigEndian, data)
	// if err != nil {
	// 	log.Println(err)
	// }
	return buf.Bytes(), nil
}
