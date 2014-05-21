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
// been buffered or if a configured time duration has elapsed.

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

	messages chan *proto.Message
	flushing chan *proto.ProduceRequest
	errors   chan error
	done     chan bool

	config AsyncProducerConfig

	nextRequest *proto.ProduceRequest

	buffered int
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
		conn:     c,
		messages: make(chan *proto.Message),
		flushing: make(chan *proto.ProduceRequest, 16),
		errors:   make(chan error, 16),
		done:     make(chan bool),
		config:   *config,
	}
	go p.run()
	go p.flush()
	return p, nil
}

func (p *AsyncProducer) Close() error {
	close(p.messages)
	<-p.done
	return p.conn.Close()
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
	defer close(p.done)
	defer close(p.errors)

	for {
		req, ok := <-p.flushing
		if !ok {
			return
		}

		b, err := proto.Encode(&proto.Request{req})
		if err != nil {
			p.errors <- err
			continue
		}

		_, err = p.conn.Write(b)
		if err != nil {
			p.errors <- err
			continue
		}
	}
}

// submitNextRequest submits the nextRequest
// for flusing.
func (p *AsyncProducer) submitNextRequest() {
	p.flushing <- p.nextRequest
	p.nextRequest = nil
	p.buffered = 0
}

// run is the main looop of the AsyncProducer
func (p *AsyncProducer) run() {
	timer := time.NewTimer(p.config.MaxBufferDuration)
	defer timer.Stop()

	// stop timer until first message is queued
	timer.Stop()

	for {
		select {
		case msg, ok := <-p.messages:
			if !ok {
				// messages channel has closed. submit any
				// pending data and exit the run loop.
				if p.nextRequest != nil {
					p.submitNextRequest()
				}
				close(p.flushing)
				return
			}

			// allocate a new ProduceRequest if one
			// doesn't exist already (i.e., start a new batch)
			if p.nextRequest == nil {
				p.nextRequest = &proto.ProduceRequest{Topic: p.config.Topic}
				timer.Reset(p.config.MaxBufferDuration)
			}

			// append the message to the next ProduceRequest
			// that will be flushed
			p.nextRequest.AddMessage(msg)

			// track the amount of data being buffered
			p.buffered += len(msg.Key)
			p.buffered += len(msg.Value)

			// if the amount of buffered data exceeds the configured
			// limit, submit the request to be flushed
			if p.buffered >= p.config.MaxBufferedBytes {
				log.Printf("buffered %d bytes, flushing", p.buffered)
				p.submitNextRequest()
			}

		case now := <-timer.C:
			// if the oldest data has existed for more than
			// MaxBufferDuration, submit it to be flushed
			log.Printf("max buffer duration at %v", now)
			if p.nextRequest != nil {
				p.submitNextRequest()
			}
		}
	}
}
