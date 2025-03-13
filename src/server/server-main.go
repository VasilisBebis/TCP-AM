package main

import (
	"fmt"
	"github.com/VasilisBebis/TCP-AM/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := server.NewServer()
	s.OpenServer()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	done := make(chan bool, 1)

	go func() {
		<-sigs
		fmt.Println()
		done <- true
	}()

	<-done

	s.CloseServer()
}
