package test

import (
	"gophq.io/gophqd"
	"log"
	"net"
)

func runServer() (net.Listener, *gophqd.Server) {
	server := &gophqd.Server{}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic("Listen: " + err.Error())
	}

	go func() {
		err := server.Serve(l)
		if err != nil {
			log.Printf("Serve: %v", err)
		}
	}()

	return l, server
}
