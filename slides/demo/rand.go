package main

import (
	"crypto/rand"
	"fmt"
	"io"
)

func Random(n int) ([]byte, error) {
	var bs = make([]byte, n)
	_, err := io.ReadFull(rand.Reader, bs)
	return bs, err
}

func main() {
	bs, err := Random(16)
	if err != nil {
		fmt.Printf("[!] rand.Reader failed: %v\n", err)
	} else {
		fmt.Printf("Data: %x\n", bs)
	}
}
