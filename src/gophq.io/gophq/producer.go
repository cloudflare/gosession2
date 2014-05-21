package gophq

import (
	"gophq.io/proto"
	"net"
	"time"
)

// Producer configuration
type ProducerConfig struct {
	// Maximum number of bytes to buffer before sending
	MaxBufferedBytes int

	// Maximum amount of time to buffer messages before sending
	MaxBufferDuration time.Duration
}

type Producer struct {
	conn   net.Conn
	errors chan error
	done   chan bool
}

func NewProducer(network, addr string, config *ProducerConfig) (*Producer, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	p := &Producer{
		conn:   c,
		errors: make(chan error),
		done:   make(chan bool),
	}
	return p, nil
}

func (p *Producer) Close() error {
	// TODO debatable whether Close should call Flush
	p.Flush()

	/// TODO signal queue sender to shut down
	close(p.done)

	return p.conn.Close()
}

// Don't mix SendMessage with QueueMessage.
func (p *Producer) SendMessage(topic string, key, value []byte) error {
	produceReq := proto.ProduceRequest{
		Topic: topic,
	}
	produceReq.AddMessage(proto.ByteEncoder(key), proto.ByteEncoder(value))

	b, err := proto.Encode(&proto.Request{&produceReq})
	if err != nil {
		return err
	}

	// TODO a socket write buffer might help here
	// TODO might want to set a write deadline

	_, err = p.conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// Don't mix QueueMessage with SendMessage.
func (p *Producer) QueueMessage(topic string, key, value []byte) {
}

func (p *Producer) Flush() {
	// TODO flush queued messages
}

// Errors returns a channel from which to read errors
// arising from QueueMessage.
// Failure to read the Errors channel will cause
// QueueMessage to block.
func (p *Producer) Errors() chan error {
	return nil
}
