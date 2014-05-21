package test

import (
	"gophq.io/gophq"
	"log"
	"os"
	"testing"
	"time"
)

func TestProducerConsumer(t *testing.T) {
	os.RemoveAll("topic.dat")

	l, server := runServer()
	defer l.Close()
	addr := l.Addr()

	config := &gophq.ConsumerConfig{
		Topic:       "topic",
		MinBytes:    280,
		MaxBytes:    280,
		MaxWaitTime: 1 * time.Second,
		FetchOffset: 0,
	}

	consumer, err := gophq.NewConsumer(addr.Network(), addr.String(), tlsConf, config)
	if err != nil {
		t.Fatalf("NewConsumer: %v", err)
	}
	defer consumer.Close()

	for {
		event := consumer.ReadEvent()
		if event.Err != nil {
			log.Printf("Err: %v", event.Err)
			break
		}
	}

	consumer.Close()
	l.Close()
	server.Terminate()
}
