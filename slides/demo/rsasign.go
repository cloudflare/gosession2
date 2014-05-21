package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
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

func main() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	checkFatal(err)

	message := []byte("Binding contractual agreement...")

	h := sha256.New()
	h.Write(message)
	digest := h.Sum(nil)
	sig, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA256, digest, nil)
	checkFatal(err)

	fmt.Printf("Signature: %x\n", sig)
	err = rsa.VerifyPSS(&priv.PublicKey, crypto.SHA256, digest, sig, nil)
	fmt.Printf("Signature OK: %v\n", err == nil)
}
