package gophq

import (
	"gophq.io/proto"
	"gophq.io/tls"
	"log"
	"net"
	"time"
)

// TODO Implement a producer that allows events to be queued,
// and flushes events if a configured number of bytes have
// been buffered (sum of key bytes and value bytes)
// or if a configured time duration has elapsed.

// Producer configuration
type AsyncProducerConfig struct {
	// Topic to produce to
	Topic string

	// Maximum number of bytes to buffer before sending
	MaxBufferedBytes int

	// Maximum amount of time to buffer messages before sending
	MaxBufferDuration time.Duration
}

// AsyncProducer
type AsyncProducer struct {
	conn net.Conn

	// TODO there are many ways to do this, but you
	// probably want at least:
	// 1. a messages channel
	// 2. an errors channel
	// 3. a "done" channel for clean shutdown

	// you might want to append incoming messages
	// to a proto.ProduceRequest
	// produceRequest.AddMessage(msg)

	config AsyncProducerConfig
}

// NewAsyncProducer creates a new producer that allows messages
// to be queued and flushed to the broker asynchronously.
func NewAsyncProducer(network, addr string, tlsConf *tls.TLSConfig,
	config *AsyncProducerConfig) (*AsyncProducer, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	// this works even if tlsConf is nil
	c = tlsConf.Client(c)

	p := &AsyncProducer{
		conn:   c,
		config: *config,
	}
	// go p.something()
	return p, nil
}

func (p *AsyncProducer) Close() error {
	// Close some channels
	// Maybe wait (read) from some channels
	// Close the connection

	return nil
}

func (p *AsyncProducer) QueueMessage(key, value []byte) {
	p.messages <- proto.NewMessage(key, value)
}

// Errors returns a channel from which to read errors
// arising from QueueMessage.
// Failure to read the Errors channel will cause
// QueueMessage to block.
func (p *AsyncProducer) Errors() chan error {
	return p.errors
}

// flush writes the next pending request to the broker.
func (p *AsyncProducer) flush() {
	// TODO you might want a function like this
}

// submitNextRequest submits the nextRequest
// for flusing.
func (p *AsyncProducer) submitNextRequest() {
	// TODO you might want a function like this
}

// run is the main looop of the AsyncProducer
func (p *AsyncProducer) run() {
	// time.NewTimer could be useful

	// for { select { } }
}

// TODO the consumer side of the server isn't quite done,
// but you could hack handleProduceRequest in
// gophqd/produce.go to just keep a list of messages
// it has received so that you can extend the
// unit tests in test/producer_test.go.
