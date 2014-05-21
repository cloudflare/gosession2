package proto

import (
	"errors"
	"fmt"
)

type KError int16

// Numeric error codes returned by the server.
const (
	NoError                 KError = 0
	Unknown                 KError = -1
	OffsetOutOfRange        KError = 1
	InvalidMessage          KError = 2
	UnknownTopicOrPartition KError = 3
	InvalidMessageSize      KError = 4
	RequestTimedOut         KError = 7
	MessageSizeTooLarge     KError = 10
)

var InsufficientData = errors.New("Insufficient data to decode packet")

var EncodingError = errors.New("Error while encoding packet")

type DecodingError struct {
	Info string
}

func (err DecodingError) Error() string {
	return fmt.Sprintf("Error while decoding packet: %s", err.Info)
}

func (err KError) Error() string {
	switch err {
	case NoError:
		return "OK"
	case OffsetOutOfRange:
		return "offset out of range"
	case InvalidMessage:
		return "invalid message (CRC mismatch)"
	case UnknownTopicOrPartition:
		return "unknown topic or partition"
	case InvalidMessageSize:
		return "invalid message size"
	case RequestTimedOut:
		return "request timed out"
	case MessageSizeTooLarge:
		return "message size too large"
	}
	return "unknown server error"
}
