package test

import (
	"gophq.io/gophq"
	"gophq.io/gophqd"
	"log"
	"net"
	"testing"
)

// TODO benchmark normal throughput

// TODO benchmark TLS throughput

// TODO look at TLS code running under perf/pprof

func TestProducer(t *testing.T) {
	server := &gophqd.Server{}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	defer l.Close()
	addr := l.Addr()

	go func() {
		err := server.Serve(l)
		if err != nil {
			log.Printf("Serve: %v", err)
		}
	}()

	config := &gophq.ProducerConfig{}
	producer, err := gophq.NewProducer(addr.Network(), addr.String(), config)
	if err != nil {
		t.Fatalf("NewProducer: %v", err)
	}
	for i := 0; i < 10; i++ {
		err = producer.SendMessage("hello", []byte("key"), []byte{byte(i)})
		if err != nil {
			t.Errorf("SendMessage: %v", err)
			break
		}
	}
	producer.Close()

	l.Close()
	server.Terminate()
}
