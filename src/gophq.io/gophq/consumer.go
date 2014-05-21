package gophq

import (
	"net"
	"time"
)

const (
	LatestOffset   = -1
	EarliestOffset = -2
)

type ConsumerConfig struct {
	// minimum fetch size
	MinBytes uint32

	// maximum fetch size
	MaxBytes uint32

	MaxWaitTime time.Duration

	OffsetValue int64
}

type Consumer struct {
	net.Conn
}

func NewConsumer(topic string, config *ConsumerConfig) *Consumer {
	net.Dial("", "")
	return nil
}

func (this *Consumer) Close() error {
	return nil
}

type ConsumerEvent struct {
	Key, Value []byte
	Offset     int64
	Err        error
}

func (this *Consumer) ReadEvent() *ConsumerEvent {
	// TODO consider setting a ReadTimeout to detect
	// a broker that doesn't adhere to MaxWaitTime
	// or detect a violation of MaxWaitTime after the fact

	return nil
}
