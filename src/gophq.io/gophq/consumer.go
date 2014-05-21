package gophq

import (
	"gophq.io/proto"
	"gophq.io/tls"
	"log"
	"net"
	"time"
)

// Consumer configuration
type ConsumerConfig struct {
	// Topic to consume from.
	Topic string

	// This is the minimum number of bytes of messages
	// that must be available to give a response.
	// If the client sets this to 0 the server will always
	// respond immediately, however if there is no new data
	// since their last request they will just get back empty
	// message sets. If this is set to 1, the server will
	// respond as soon as the topic has at least 1 byte of
	// data or the specified timeout occurs.
	// By setting higher values in combination with the
	// imeout the consumer can tune for throughput and trade a
	// little additional latency for reading only large chunks of data.
	MinBytes int32

	// maximum fetch size
	MaxBytes int32

	// The max wait time is the maximum amount of time in
	// milliseconds to block waiting if insufficient data
	// is available at the time the request is issued.
	MaxWaitTime time.Duration

	// The offset to begin this fetch from.
	FetchOffset int64
}

// Consumer fetches data from the broker. Consumer
// embeds net.Conn, so it "inherits" its Close function.
type Consumer struct {
	net.Conn
}

// NewConsumer creates a new connection
// to the broker for the specified topic.
func NewConsumer(network, addr string, tlsConf *tls.TLSConfig, config *ConsumerConfig) (*Consumer, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	c = tlsConf.Client(c)

	fetchReq := &proto.FetchRequest{
		Topic:       config.Topic,
		MinBytes:    config.MinBytes,
		MaxBytes:    config.MaxBytes,
		MaxWaitTime: config.MaxWaitTime,
		FetchOffset: config.FetchOffset,
	}

	b, err := proto.Encode(&proto.Request{fetchReq})
	if err != nil {
		return nil, err
	}

	// TODO a socket write buffer might help here
	// TODO might want to set a write deadline

	_, err = c.Write(b)
	if err != nil {
		return nil, err
	}

	return &Consumer{c}, nil
}

// ConsumerEvent contains the next
// Key/Value and the offset from which it was fetched,
// or an error in the Err field.
type ConsumerEvent struct {
	Key, Value []byte
	Offset     int64
	Err        error
}

// ReadEvent reads the next ConsumerEvent
// from the consumer connection.
func (this *Consumer) ReadEvent() *ConsumerEvent {
	b, err := proto.ReadRequestOrResponse(this)
	if err != nil {
		return &ConsumerEvent{Err: err}
	}

	var fetchResp proto.FetchResponse
	err = proto.Decode(b, &fetchResp)
	if err != nil {
		return &ConsumerEvent{Err: err}
	}

	if fetchResp.Err != proto.NoError {
		return &ConsumerEvent{Err: fetchResp.Err}
	}

	log.Printf("%+v", fetchResp)

	return nil
}
