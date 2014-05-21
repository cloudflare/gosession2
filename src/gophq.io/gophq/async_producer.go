package gophq

import (
	"gophq.io/proto"
	"log"
	"net"
	"time"
)

// Producer configuration
type AsyncProducerConfig struct {
	Topic string

	// Maximum number of bytes to buffer before sending
	MaxBufferedBytes int

	// Maximum amount of time to buffer messages before sending
	MaxBufferDuration time.Duration
}

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

func NewAsyncProducer(network, addr string, config *AsyncProducerConfig) (*AsyncProducer, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
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

func (p *AsyncProducer) flushNextRequest() {
	p.flushing <- p.nextRequest
	p.nextRequest = nil
	p.buffered = 0
}

func (p *AsyncProducer) run() {
	timer := time.NewTimer(p.config.MaxBufferDuration)
	defer timer.Stop()

	// stop timer until first message is queued
	timer.Stop()

	for {
		select {
		case msg, ok := <-p.messages:
			if !ok {
				if p.nextRequest != nil {
					p.flushNextRequest()
				}
				close(p.flushing)
				return
			}

			if p.nextRequest == nil {
				p.nextRequest = &proto.ProduceRequest{Topic: p.config.Topic}
				timer.Reset(p.config.MaxBufferDuration)
			}

			p.nextRequest.AddMessage(msg)
			p.buffered += len(msg.Key)
			p.buffered += len(msg.Value)

			if p.buffered >= p.config.MaxBufferedBytes {
				log.Printf("buffered %d bytes, flushing", p.buffered)
				p.flushNextRequest()
			}

		case now := <-timer.C:
			log.Printf("max buffer duration at %v", now)
			if p.nextRequest != nil {
				p.flushNextRequest()
			}
		}
	}
}
