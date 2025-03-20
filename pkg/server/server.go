// Package that implements the server's side logic of the program
package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/client"
	"log"
	"net"
	"strconv"
)

// Header consists of the header fields used
// in the server's response message to the client
type Header struct {
	Response_code  byte   // indicates if the client's query was successful
	Length         byte   // length of the data portion of the message in bytes (padding EXCLUDED)
	Transaction_id []byte // unique identifier for the message (query & response)
}

// Message represents the full response message
type Message struct {
	Header  Header
	Result  []byte // message data serialized as a byte array
	Padding []byte // used to make the message 32-bit aligned (empty if not needed)
}

// Generates a new Message object that includes the given result
func NewMessage(response_code byte, transaction_id []byte, result any) *Message {

	bin_result, err := SerializeResult(result)
	if err != nil {
		return nil
	}
	length := len(bin_result)
	padding_size := (4 - length%4) % 4
	padding := make([]byte, padding_size)
	header := Header{Response_code: response_code, Length: byte(length), Transaction_id: transaction_id}
	message := Message{Header: header, Result: bin_result, Padding: padding}
	return &message
}

// SerializeResult creates a byte array of the given result
func SerializeResult(result any) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch v := result.(type) {
	case int32:
		binary.Write(buf, binary.BigEndian, v)
	case float32:
		binary.Write(buf, binary.BigEndian, v)
	case []int32:
		for _, i := range v {
			binary.Write(buf, binary.BigEndian, i)
		}
	default:
		return nil, fmt.Errorf("Type %T is not supported!", v)
	}
	return buf.Bytes(), nil
}

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

	message := client.DeserializeMessage(buf)
	header := message.Header
	op_code := header.Op_code
	// length := header.Length
	// transaction_id := header.Transaction_id

	result, err := CalculateResult(op_code, message.Data)
	_ = result

	fmt.Printf("%s", buf)
}

// CalculateResult calculates the result of the given operation and returns it (or returns error if there is any)
func CalculateResult(op_code byte, data []byte) (any, error) {
	var data_type string
	if op_code == 0 {
		data_type = "int8"
		data, err := client.DeserializeData(data, data_type)
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
		data, err := client.DeserializeData(data, data_type)
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
		data, err := client.DeserializeData(data, data_type)
		if err != nil {
			log.Println(err)
		}

		raw_data, ok := data.([]uint16)
		if ok {
			set_1 := raw_data[:len(raw_data)/2-1]
			set_2 := raw_data[len(raw_data):]
			result := make([]int32, len(raw_data))
			for i := 0; i < len(raw_data); i++ {
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
	// s.Conn.Close()
	err := s.Listener.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nServer Closing")
}
