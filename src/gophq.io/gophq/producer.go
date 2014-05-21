package gophq

import (
	"gophq.io/proto"
	"net"
)

// TODO make Producer safe to use from multiple goroutines

type Producer struct {
	net.Conn
}

func NewProducer(network, addr string) (*Producer, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &Producer{c}, nil
}

func (p *Producer) SendMessage(topic string, key, value []byte) error {
	produceReq := proto.ProduceRequest{
		Topic: topic,
	}
	produceReq.AddMessage(proto.NewMessage(key, value))

	b, err := proto.Encode(&proto.Request{&produceReq})
	if err != nil {
		return err
	}

	// TODO a socket write buffer might help here
	// TODO might want to set a write deadline

	_, err = p.Conn.Write(b)
	if err != nil {
		return err
	}
	return nil
}
