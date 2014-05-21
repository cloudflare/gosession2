package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

func Random(n int) ([]byte, error) {
	var bs = make([]byte, n)
	_, err := io.ReadFull(rand.Reader, bs)
	return bs, err
}

func Hash(bs []byte) []byte {
	h := sha256.New()
	h.Write(bs)
	return h.Sum(nil)
}

func main() {
	bs, err := Random(16)
	if err != nil {
		fmt.Printf("[!] rand.Reader failed: %v\n", err)
		return
	}
	digest := Hash(bs)
	fmt.Printf("Data: %x\nDigest: %x\n", bs, digest)
}
