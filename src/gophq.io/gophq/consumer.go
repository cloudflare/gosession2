package gophq

import (
	"gophq.io/proto"
	"log"
	"net"
	"time"
)

const (
	LatestOffset   = -1
	EarliestOffset = -2
)

type ConsumerConfig struct {
	// minimum fetch size
	MinBytes int32

	// maximum fetch size
	MaxBytes int32

	MaxWaitTime time.Duration

	FetchOffset int64
}

type Consumer struct {
	net.Conn
}

func NewConsumer(network, addr string, topic string, config *ConsumerConfig) (*Consumer, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	fetchReq := &proto.FetchRequest{
		Topic:       topic,
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

type ConsumerEvent struct {
	Key, Value []byte
	Offset     int64
	Err        error
}

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
