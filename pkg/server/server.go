// Package that implements the server's side logic of the program
package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	// "sync"
)

// Header consists of the header fields used
// in the server's response message to the client
type Header struct {
	Response_code  byte    // indicates if the client's query was successful
	Length         byte    // length of the data portion of the message in bytes (padding EXCLUDED)
	Transaction_id [2]byte // unique identifier for the message (query & response)
}

// Message represents the full response message
type Message struct {
	Header  Header
	Data    []byte // message data serialized as a byte array
	Padding []byte // used to make the message 32-bit aligned (empty if not needed)
}

// Default Port used if the port is not specified
const Def_Port = "12345"

type Server struct {
	Listener net.Listener
	// Conn     net.Conn
	Port  string
	Close bool
}

// NewServer creates a server that will be able to listen on the given port.
// By default the server is not opened. To open the server use [OpenServer]
func NewServer() *Server {
	server := Server{Listener: nil, Port: Def_Port, Close: false}
	return &server
}

// ChangePort changes the TCP port that the server will listen on.
func (s *Server) ChangePort(port uint64) {
	s.Port = strconv.FormatUint(port, 10)
}

// OpenServer opens the listening ability on the TCP port of the given server on all available unicast
// and anycast IP addresses of the local system.
func (s *Server) OpenServer() {
	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server listening at port " + s.Port)
	s.Listener = listener

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if s.Close {
					return
				}
				fmt.Println("Connection Error:", err)
			}
			go handleClient(conn)
		}
	}()
}

func handleClient(c net.Conn) {
	defer c.Close()
	fmt.Println("Client: ", c.RemoteAddr(), " connected")

}

func (s *Server) CloseServer() {
	s.Close = true
	// s.Conn.Close()
	err := s.Listener.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nServer Closing")
}

func HelloServer() {
	fmt.Println("Hello from Server")
}
