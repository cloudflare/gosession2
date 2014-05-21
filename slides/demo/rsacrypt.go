package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func fatalf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}

func checkFatal(err error) {
	if err == nil {
		return
	}
	fatalf("%v\n", err)
}

func Random(n int) ([]byte, error) {
	var bs = make([]byte, n)
	_, err := io.ReadFull(rand.Reader, bs)
	return bs, err
}

func main() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	checkFatal(err)

	// 16-byte AES key + 32-byte HMAC-SHA-256 key
	sessionKey, err := Random(48)
	checkFatal(err)

	h := sha256.New()
	out, err := rsa.EncryptOAEP(h, rand.Reader, &priv.PublicKey, sessionKey, nil)
	checkFatal(err)

	key, err := rsa.DecryptOAEP(h, rand.Reader, priv, out, nil)
	checkFatal(err)

	if !bytes.Equal(key, sessionKey) {
		fatalf("Decrypted key doesn't match original key!\n")
	}
	fmt.Println("OK")
}
