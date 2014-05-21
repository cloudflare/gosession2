package gophq

import (
	"gophq.io/proto"
	"gophq.io/tls"
	"log"
	"net"
)

type Producer struct {
	net.Conn
}

// NewProducer establishes a new producer connection
// to the broker.
func NewProducer(network, addr string, tlsConf *tls.TLSConfig) (*Producer, error) {
	c, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	// this works even if tlsConf is nil
	c = tlsConf.Client(c)

	return &Producer{c}, nil
}

// SendMessage sends a ProduceRequest with a MessageSet containing
// a single Message (Key/Value) to the broker.
func (p *Producer) SendMessage(topic string, key, value []byte) error {
	produceReq := &proto.ProduceRequest{
		Topic: topic,
	}
	produceReq.AddMessage(proto.NewMessage(key, value))

	b, err := proto.Encode(&proto.Request{produceReq})
	if err != nil {
		return err
	}

	// TODO a socket write buffer might help here
	// TODO might want to set a write deadline

	_, err = p.Conn.Write(b)
	if err != nil {
		return err
	}

	b, err = proto.ReadRequestOrResponse(p)
	if err != nil {
		return err
	}

	var produceResp proto.ProduceResponse
	err = proto.Decode(b, &produceResp)
	if err != nil {
		return err
	}

	log.Printf("%+v", produceResp)

	// Offset in ProduceResponse is not returned for now

	return nil
}
