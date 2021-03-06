package proto

import (
	"encoding/binary"
	"hash/crc32"
)

// crc32Field implements the pushEncoder and pushDecoder interfaces for calculating CRC32s.
type crc32Field struct {
	startOffset int
	crc         uint32
}

func (c *crc32Field) saveOffset(in int) {
	c.startOffset = in
}

func (c *crc32Field) reserveLength() int {
	return 4
}

// TODO this can be much improved
// first step: crc32.MakeTable(crc32.Castagnoli)

func (c *crc32Field) run(curOffset int, buf []byte) error {
	c.crc = crc32.ChecksumIEEE(buf[c.startOffset+4 : curOffset])
	binary.BigEndian.PutUint32(buf[c.startOffset:], c.crc)
	return nil
}

func (c *crc32Field) check(curOffset int, buf []byte) error {
	c.crc = crc32.ChecksumIEEE(buf[c.startOffset+4 : curOffset])

	if c.crc != binary.BigEndian.Uint32(buf[c.startOffset:]) {
		return DecodingError{Info: "CRC32 didn't match"}
	}

	return nil
}
