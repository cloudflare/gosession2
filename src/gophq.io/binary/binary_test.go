package gophqbinary

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"
)

func TestBigEndian(t *testing.T) {
	buffer := make([]byte, 14)
	PutBigEndianUint16(buffer, 0, 0x0102)
	PutBigEndianUint32(buffer, 2, 0x03040506)
	PutBigEndianUint64(buffer, 6, 0x0708090a0b0c0d0e)

	expected := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e}
	if !bytes.Equal(expected, buffer) {
		t.Fatal("incorrect result")
	}

	x := GetBigEndianUint16(buffer, 0)
	if x != 0x0102 {
		t.Fatalf("GetBigEndianUint16: x=%x", x)
	}

	y := GetBigEndianUint32(buffer, 2)
	if y != 0x03040506 {
		t.Fatalf("GetBigEndianUint32: y=%x", y)
	}

	z := GetBigEndianUint64(buffer, 6)
	if z != 0x0708090a0b0c0d0e {
		t.Fatalf("GetBigEndianUint64: z=%x", z)
	}
}

func BenchmarkGetBE16(b *testing.B) {
	buffer := []byte{0x00, 0x01, 0x02}
	for k := 0; k < b.N; k++ {
		x := GetBigEndianUint16(buffer, 1)
		_ = x
	}
}

func BenchmarkGoGetBE16(b *testing.B) {
	buffer := []byte{0x00, 0x01, 0x02}
	for k := 0; k < b.N; k++ {
		x := binary.BigEndian.Uint16(buffer[1:])
		_ = x
	}
}

func BenchmarkGetBE32(b *testing.B) {
	buffer := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	for k := 0; k < b.N; k++ {
		x := GetBigEndianUint32(buffer, 1)
		_ = x
	}
}

func BenchmarkGoGetBE32(b *testing.B) {
	buffer := []byte{0x00, 0x01, 0x02, 0x03, 0x04}
	for k := 0; k < b.N; k++ {
		x := binary.BigEndian.Uint32(buffer[1:])
		_ = x
	}
}

func BenchmarkGetBE64(b *testing.B) {
	buffer := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	for k := 0; k < b.N; k++ {
		x := GetBigEndianUint64(buffer, 1)
		_ = x
	}
}

func BenchmarkGoGetBE64(b *testing.B) {
	buffer := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	for k := 0; k < b.N; k++ {
		x := binary.BigEndian.Uint64(buffer[1:])
		_ = x
	}
}

func TestSetBigEndianUint16(t *testing.T) {
	const v = 0x0123
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, v)
	expected := binary.LittleEndian.Uint16(data)
	var actual uint16
	SetBigEndianUint16(&actual, v)
	if expected != actual {
		t.Fatalf("expected %x, got %x", expected, actual)
	}
}

func TestSetBigEndianUint32(t *testing.T) {
	const v = 0x01234567
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, v)
	expected := binary.LittleEndian.Uint32(data)
	var actual uint32
	SetBigEndianUint32(&actual, v)
	if expected != actual {
		t.Fatalf("expected %x, got %x", expected, actual)
	}
}

func TestSetBigEndianUint64(t *testing.T) {
	const v = 0x0123456789abcdef
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, v)
	expected := binary.LittleEndian.Uint64(data)
	var actual uint64
	SetBigEndianUint64(&actual, v)
	if expected != actual {
		t.Fatalf("expected %x, got %x", expected, actual)
	}
}

func BenchmarkBigEndianPutUint64(b *testing.B) {
	const v = 0x0123456789abcdef
	var dst uint64
	data := (*(*[8]byte)(unsafe.Pointer(&dst)))[:]
	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint64(data, v)
	}
}

func BenchmarkSetBigEndianUint64(b *testing.B) {
	const v = 0x0123456789abcdef
	var dst uint64
	for i := 0; i < b.N; i++ {
		SetBigEndianUint64(&dst, v)
	}
}
