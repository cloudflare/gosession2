package test

import (
	"flag"
	"gophq.io/gophqd"
	"gophq.io/tls"
	"log"
	"net"
)

var testTLS = flag.Bool("tls", true, "test with TLS")

var tlsConf *tls.TLSConfig

func init() {
	flag.Parse()
	if *testTLS {
		tlsConf = tls.SelfSignedTLSConfig()
	}
}

func runServer() (net.Listener, *gophqd.Server) {
	server := &gophqd.Server{
		TLSConfig: tlsConf,
	}

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
