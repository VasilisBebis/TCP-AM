// Package that implements the client's side logic of the program
package client

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/messages"
	"log"
	"net"
)

type Client struct {
	Conn net.Conn // client's connector
}

// NewClient creates an object of type [Client] and return it
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

// SendQuery sends the given query message to the server
func (c *Client) SendQuery(q messages.Query) {
	ser_query := q.SerializeQuery()
	_, err := c.Conn.Write(ser_query)
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
