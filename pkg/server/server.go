// Package that implements the server's side logic of the program
package server

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/messages"
	"log"
	"net"
	"strconv"
)

// Default Port used if the port is not specified
const Def_Port = "12345"

type Server struct {
	Listener net.Listener
	Port     string
	Close    bool
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
	fmt.Println("Client: ", c.RemoteAddr(), " connected!")

	buf := make([]byte, 256)
	_, err := c.Read(buf)
	if err != nil {
		log.Println(err)
	}

	message := messages.DeserializeQuery(buf)
	header := message.Header
	op_code := header.Op_code
	transaction_id := header.Transaction_id

	result, err := CalculateResult(op_code, message.Data)
	var r_ptr *messages.Response
	if err != nil {
		r_ptr = messages.NewResponse(0, transaction_id, []byte{})
	}
	//TODO: create appropriate error messages
	r_ptr = messages.NewResponse(0, transaction_id, result)
	r := *r_ptr
	ser_response := r.SerializeResponse()
	_, err = c.Write(ser_response)
	if err != nil {
		log.Println(err)
	}
}

// CalculateResult calculates the result of the given operation and returns it (or returns error if there is any)
func CalculateResult(op_code byte, data []byte) (any, error) {
	var data_type string
	if op_code == 0 {
		data_type = "int8"
		data, err := messages.DeserializeData(data, data_type)
		if err != nil {
			log.Println(err)
		}

		raw_data, ok := data.([]int8)
		if ok {
			var product int32
			product = 1
			for _, v := range raw_data {
				if v < -5 || v > 5 {
					return nil, fmt.Errorf("Given number out of range!")
				}
				product *= int32(v)
			}
			return product, nil
		}
	} else if op_code == 1 {
		data_type = "uint8"
		data, err := messages.DeserializeData(data, data_type)
		if err != nil {
			log.Println(err)
		}

		raw_data, ok := data.([]uint8)
		if ok {
			var sum uint16
			sum = 0
			count := 0
			for _, v := range raw_data {
				if v > 200 {
					return nil, fmt.Errorf("Given number out of range!")
				}
				sum += uint16(v)
				count++
			}
			avg := float32(sum) / float32(count)
			return avg, nil
		}

	} else if op_code == 2 {
		data_type = "uint16"
		data, err := messages.DeserializeData(data, data_type)
		if err != nil {
			log.Println(err)
		}

		raw_data, ok := data.([]uint16)
		if ok {
			set_1 := raw_data[:len(raw_data)/2-1]
			set_2 := raw_data[len(raw_data):]
			result := make([]int32, len(raw_data))
			for i := range len(raw_data) {
				if set_1[i] > 60_000 || set_2[i] > 60_000 {
					return nil, fmt.Errorf("Given number out of range!")
				}
				result[i] = int32(set_1[i] - set_2[i])
			}
			return result, nil
		}
	}
	return nil, fmt.Errorf("Invalid operation!")
}

func (s *Server) CloseServer() {
	s.Close = true
	err := s.Listener.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nServer Closing")
}
