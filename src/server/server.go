// Package that implements the server's side logic of the program
package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

// Default Port used if the port is not specified
const Def_Port = "12345"

type Server struct {
	Listener net.Listener
	Conn     net.Conn
	Port     string
	Close    bool
}

// NewServer creates a server that will be able to listen on the given port.
// By default the server is not opened. To open the server use [OpenServer]
func NewServer() *Server {
	server := Server{Listener: nil, Conn: nil, Port: Def_Port, Close: false}
	return &server

}

// ChangePort changes the TCP port that the server will listen on.
func (s *Server) ChangePort(port uint64) {
	s.Port = strconv.FormatUint(port, 10)
}

// OpenServer opens the listening ability on the TCP port of the given server on all available unicast
// and anycast IP addresses of the local system.
// It also initiates a connection.
func (s *Server) OpenServer() {
	listener, err := net.Listen("tcp", ":"+Def_Port)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	s.Conn = conn
	s.Listener = listener
	fmt.Println("Server listening at port " + Def_Port)
}

func CloseServer(s *Server) {
	(*s).Conn.Close()
	(*s).Listener.Close()
}

func HelloServer() {
	fmt.Println("Hello from Server")
}
