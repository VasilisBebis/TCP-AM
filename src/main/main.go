package main

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/src/client"
	"github.com/VasilisBebis/TCP-AM/src/server"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	serv := server.NewServer()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serv.OpenServer()
	}()
	time.Sleep(time.Second)
	wg.Add(1)
	defer serv.CloseServer()

	go func() {
		defer wg.Done()
		cl := client.NewClient()
		cl.ConnectTo(*serv)
		cl.Conn.Write([]byte("Hello Server"))
		buf := make([]byte, 1024)
		_, err := serv.Conn.Read(buf)
		_ = err
		fmt.Printf("%s", buf)
	}()

	wg.Wait()

}
