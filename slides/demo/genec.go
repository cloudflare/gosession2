package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fatalf("%v\n", err)
	}

	out, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		fatalf("%v\n", err)
	}

	var p = &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: out,
	}
	fmt.Printf("%s\n", pem.EncodeToMemory(p))
}
