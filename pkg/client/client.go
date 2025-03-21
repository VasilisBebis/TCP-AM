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

// DeserializeMessage unpacks a byte stream to an object of type [Message]
func DeserializeMessage(byte_stream []byte) *Message {
	op_code := byte_stream[0]
	length := byte_stream[1]
	transaction_id := byte_stream[2:3]
	data := byte_stream[4:(4 + length)]
	h := Header{Op_code: op_code, Length: length, Transaction_id: transaction_id}
	m := Message{Header: h, Data: data}
	return &m
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
	fmt.Println("Connected to server ", server_ip, ":", server_port)
	c.Conn = conn
}

// SendMessage sends the given message to the server
func (c *Client) SendMessage(m Message) {
	ser_message := m.SerializeMessage()
	_, err := c.Conn.Write(ser_message)
	if err != nil {
		log.Println(err)
	}
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
	return buf.Bytes(), nil
}

// DeserializeData returns the data byte stream in its original format
func DeserializeData(data []byte, data_type string) (any, error) {
	var des_data any
	if data_type == "int8" {
		des_data = make([]int8, len(data))
	} else if data_type == "uint8" {
		des_data = make([]uint8, len(data))
	} else if data_type == "uint16" {
		des_data = make([]uint16, len(data)/2)
	} else {
		return nil, fmt.Errorf("Type %s is not supported", data_type)
	}

	_, err := binary.Decode(data, binary.BigEndian, des_data)
	if err != nil {
		log.Println(err)
	}
	return des_data, nil
}
