package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
)

func Random(n int) ([]byte, error) {
	var bs = make([]byte, n)
	_, err := io.ReadFull(rand.Reader, bs)
	return bs, err
}

func HMAC(in, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(in)
	return h.Sum(nil)
}

func NewKey(h hash.Hash) ([]byte, error) {
	return Random(h.Size())
}

func main() {
	message := []byte("Hello, world.")
	key, err := NewKey(sha256.New())
	if err != nil {
		fmt.Printf("Random failed: %v\n", err)
		return
	}
	tag := HMAC(message, key)
	fmt.Printf("  Tag: %x\n", tag)

	message2 := []byte("Hello, world-")
	tag2 := HMAC(message2, key)
	fmt.Printf("Tag 2: %x\n", tag2)
	fmt.Printf("tag1 == tag2: %v\n", hmac.Equal(tag, tag2))

	key, err = NewKey(sha256.New())
	if err != nil {
		fmt.Printf("Random failed: %v\n", err)
		return
	}
	tag2 = HMAC(message, key)
	fmt.Printf("Tag 2: %x\n", tag2)
	fmt.Printf("tag1 == tag2: %v\n", hmac.Equal(tag, tag2))
}
