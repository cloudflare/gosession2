package proto

// https://github.com/Shopify/sarama/blob/master/packet_decoder.go

type packetDecoder interface {
	getInt8() (int8, error)
	getInt16() (int16, error)
	getInt32() (int32, error)
	getInt64() (int64, error)
	getArrayLength() (int, error)
	getBytes() ([]byte, error)
	getString() (string, error)
	getInt32Array() ([]int32, error)
	getInt64Array() ([]int64, error)
	remaining() int
	getSubset(length int) (packetDecoder, error)
	push(in pushDecoder) error
	pop() error
}

type pushDecoder interface {
	saveOffset(in int)
	reserveLength() int
	check(curOffset int, buf []byte) error
}
