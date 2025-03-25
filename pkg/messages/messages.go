// Package declares client and server messages and implements basic properties of them
package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand/v2"
)

type Op_code = byte

const (
	A Op_code = iota
	B
	C
)

// ClientHeader consists of the header fields used
// in the client's query message to the server
type ClientHeader struct {
	Op_code        Op_code // operation code (0 -> A, 1 -> B, 2 -> C)
	Length         byte    // length of the data portion of the message in bytes (padding EXCLUDED)
	Transaction_id []byte  // unique identifier for the client-server interaction
}

// Query represents the full client's query message
type Query struct {
	Header  ClientHeader
	Data    []byte // query data serialized as a byte array
	Padding []byte // used to make the message 32-bit aligned (empty if not needed)
}

// Generates a new [Query] object that includes the given data
func NewQuery(op byte, data any) *Query {
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
	header := ClientHeader{Op_code: op, Length: byte(length), Transaction_id: bin_tr_id.Bytes()}
	message := Query{Header: header, Data: bin_data, Padding: padding}
	return &message
}

// SerializeQuery packs an object of type [Query] to a byte stream (big endian)
func (q *Query) SerializeQuery() []byte {
	var message_bytes []byte

	message_bytes = append(message_bytes, []byte{q.Header.Op_code, q.Header.Length}...)
	message_bytes = append(message_bytes, q.Header.Transaction_id[:]...)
	message_bytes = append(message_bytes, q.Data...)
	if len(q.Padding) != 0 {
		message_bytes = append(message_bytes, q.Padding...)
	}
	return message_bytes
}

// DeserializeQuery unpacks a byte stream to an object of type [Query]
func DeserializeQuery(byte_stream []byte) *Query {
	op_code := byte_stream[0]
	length := byte_stream[1]
	transaction_id := byte_stream[2:4]
	data := byte_stream[4:(4 + length)]
	h := ClientHeader{Op_code: op_code, Length: length, Transaction_id: transaction_id}
	q := Query{Header: h, Data: data}
	return &q
}

// SerializeData creates a byte array of the given [Query] data
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

type Response_code = byte

const (
	Ok Response_code = iota
	NumberOutOfBounds
	ListLengthOutOfBounds
)

var response_message = map[Response_code]string{
	Ok:                    "Query was successful!",
	NumberOutOfBounds:     "One or more of the given numbers is out of bounds!",
	ListLengthOutOfBounds: "Given list(s) have a length that is not permitted!",
}

// ServerHeader consists of the header fields used
// in the server's response message to the client
type ServerHeader struct {
	Response_code  Response_code // indicates if the client's query was successful
	Length         byte          // length of the data portion of the message in bytes (padding EXCLUDED)
	Transaction_id []byte        // unique identifier for the client-server interaction
}

// Response represents the full server's response message
type Response struct {
	Header  ServerHeader
	Result  []byte // response data serialized as a byte array
	Padding []byte // used to make the message 32-bit aligned (empty if not needed)
}

// Generates a new [Response] object that includes the given result
func NewResponse(response_code byte, transaction_id []byte, result any) *Response {
	bin_result, err := SerializeResult(result)
	if err != nil {
		log.Println(err)
		return nil
	}
	length := len(bin_result)
	padding_size := (4 - length%4) % 4
	padding := make([]byte, padding_size)
	header := ServerHeader{Response_code: response_code, Length: byte(length), Transaction_id: transaction_id}
	message := Response{Header: header, Result: bin_result, Padding: padding}
	return &message
}

// SerializeResponse packs an object of type [Response] to a byte stream (big endian)
func (r *Response) SerializeResponse() []byte {
	var message_bytes []byte

	message_bytes = append(message_bytes, []byte{r.Header.Response_code, r.Header.Length}...)
	message_bytes = append(message_bytes, r.Header.Transaction_id[:]...)
	message_bytes = append(message_bytes, r.Result...)
	if len(r.Padding) != 0 {
		message_bytes = append(message_bytes, r.Padding...)
	}
	return message_bytes
}

// DeserializeResponse unpacks a byte stream to an object of type [Response]
func DeserializeResponse(byte_stream []byte) *Response {
	respones_code := byte_stream[0]
	length := byte_stream[1]
	transaction_id := byte_stream[2:4]
	result := byte_stream[4:(4 + length)]
	h := ServerHeader{Response_code: respones_code, Length: length, Transaction_id: transaction_id}
	r := Response{Header: h, Result: result}
	return &r
}

// SerializeResult creates a byte array of the given result
func SerializeResult(result any) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch v := result.(type) {
	case int32:
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			log.Println(err)
		}
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

// DeserializeResult returns the result byte stream in its original format
func DeserializeResult(result []byte, op_code byte) any {
	if op_code == 0 {
		var raw int32

		buf := bytes.NewReader(result)
		err := binary.Read(buf, binary.BigEndian, &raw)
		if err != nil {
			log.Println(err)
		}
		return raw

	} else if op_code == 1 {
		//TODO: implement this

	} else if op_code == 2 {
		//TODO: implement this

	}
	return nil
}
