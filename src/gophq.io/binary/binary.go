package gophqbinary

func GetBigEndianUint16(b []byte, offset int) uint16

func PutBigEndianUint16(b []byte, offset int, v uint16)

func GetBigEndianUint32(b []byte, offset int) uint32

func PutBigEndianUint32(b []byte, offset int, v uint32)

func GetBigEndianUint64(b []byte, offset int) uint64

func PutBigEndianUint64(b []byte, offset int, v uint64)

func GetBigEndianInt64(b []byte, offset int) int64 {
	return int64(GetBigEndianUint64(b, offset))
}

func PutBigEndianInt64(b []byte, offset int, v int64) {
	PutBigEndianUint64(b, offset, uint64(v))
}

func SetBigEndianUint16(dst *uint16, v uint16)

func SetBigEndianUint32(dst *uint32, v uint32)

func SetBigEndianUint64(dst *uint64, v uint64)
