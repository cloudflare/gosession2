package test

import (
	"gophq.io/gophq"
	"log"
	"testing"
	"time"
)

// These tests don't consume the events produced,
// so they are really just basic sanity checks.

func TestProducer(t *testing.T) {
	l, server := runServer()
	defer l.Close()
	addr := l.Addr()

	producer, err := gophq.NewProducer(addr.Network(), addr.String())
	if err != nil {
		t.Fatalf("NewProducer: %v", err)
	}
	for i := 0; i < 10; i++ {
		err = producer.SendMessage("topic", []byte("key"), []byte{byte(i)})
		if err != nil {
			t.Errorf("SendMessage: %v", err)
			break
		}
	}
	producer.Close()

	l.Close()
	server.Terminate()
}

func testAsyncProducer(t *testing.T, config *gophq.AsyncProducerConfig) {
	l, server := runServer()
	defer l.Close()
	addr := l.Addr()

	producer, err := gophq.NewAsyncProducer(addr.Network(), addr.String(), config)
	if err != nil {
		t.Fatalf("NewAsyncProducer: %v", err)
	}

	go func() {
		errors := producer.Errors()
		for {
			err, ok := <-errors
			if !ok {
				return
			}
			log.Printf("Error: %v", err)
		}
	}()

	for i := 0; i < 10; i++ {
		producer.QueueMessage([]byte("key"), []byte{byte(i)})
	}

	producer.Close()
	l.Close()
	server.Terminate()
}

func TestAsyncProducerBigBuffer(t *testing.T) {
	config := &gophq.AsyncProducerConfig{
		Topic:             "topic",
		MaxBufferedBytes:  100000,
		MaxBufferDuration: 1 * time.Hour,
	}
	testAsyncProducer(t, config)
}

func TestAsyncProducerNoBuffer(t *testing.T) {
	config := &gophq.AsyncProducerConfig{
		Topic:             "topic",
		MaxBufferedBytes:  0,
		MaxBufferDuration: 1 * time.Hour,
	}
	testAsyncProducer(t, config)
}

func TestAsyncProducerShortBufferDuration(t *testing.T) {
	config := &gophq.AsyncProducerConfig{
		Topic:             "topic",
		MaxBufferedBytes:  1 << 30,
		MaxBufferDuration: 1 * time.Nanosecond,
	}
	testAsyncProducer(t, config)
}
