package proto

// https://github.com/Shopify/sarama/blob/master/packet_encoder.go

type packetEncoder interface {
	putInt8(in int8)
	putInt16(in int16)
	putInt32(in int32)
	putInt64(in int64)
	putArrayLength(in int) error
	putBytes(in []byte) error
	putRawBytes(in []byte) error
	putString(in string) error
	putInt32Array(in []int32) error
	putInt64Array(in []int64) error
	push(in pushEncoder)
	pop() error
}

type pushEncoder interface {
	saveOffset(in int)
	reserveLength() int
	run(curOffset int, buf []byte) error
}
