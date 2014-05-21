package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func fatalf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}

func main() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fatalf("%v\n", err)
	}

	out := x509.MarshalPKCS1PrivateKey(priv)
	var p = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: out,
	}
	fmt.Printf("%s\n", pem.EncodeToMemory(p))
}
