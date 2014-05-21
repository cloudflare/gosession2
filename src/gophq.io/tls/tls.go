package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"time"
)

type TLSConfig struct {
	CA   []byte
	Cert []byte
	Key  []byte
}

// Certificate creates a X.509 Key Pair from
// Cert and Key bytes and prepares the Certificate.
func (this *TLSConfig) Certificate() (tlsCert tls.Certificate, err error) {
	return
}

var errNotPEMData = errors.New("not PEM data")

// CACertificate decodes the PEM block in the CA bytes
// and parses a X.509 certificate from the result.
func (this *TLSConfig) CACertificate() (*x509.Certificate, error) {
	return nil, nil
}

// Client takes a net.Conn and returns a
// crypto/tls.Conn for use on the client side of
// a TLS connection.
func (this *TLSConfig) Client(c net.Conn) net.Conn {
	if this == nil {
		log.Printf("client not using TLS")
		return c
	}

	log.Printf("client using TLS")

	return nil
}

// Server takes a net.Conn and returns a
// crypto/tls.Conn for use on the server side of
// a TLS connection.
func (this *TLSConfig) Server(c net.Conn) net.Conn {
	if this == nil {
		log.Printf("server not using TLS")
		return c
	}

	log.Printf("server using TLS")

	return nil
}

func NewTLSConfig(ca, cert, key string) *TLSConfig {
	if ca == "" || cert == "" || key == "" {
		return nil
	}

	config := &TLSConfig{}
	var err error

	config.CA, err = ioutil.ReadFile(ca)
	if err != nil {
		panic(err)
	}

	config.Cert, err = ioutil.ReadFile(cert)
	if err != nil {
		panic(err)
	}

	config.Key, err = ioutil.ReadFile(key)
	if err != nil {
		panic(err)
	}

	return config
}

// Cipher suites that are made fast by the AES-NI instructions
var fastCipherSuites = []uint16{
	tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
}

// Lots of hints in $GOROOT/src/pkg/crypto/tls/generate_cert.go

// SelfSignedTLSConfig generates a self-signed
// CA, certificate and key. Very useful for unit testing.
func SelfSignedTLSConfig() *TLSConfig {
	const rsaBits = 2048
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		panic(err)
	}

	endOfTime := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),

		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},

		NotBefore: time.Now(),
		NotAfter:  endOfTime,

		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},

		BasicConstraintsValid: true,
	}

	template.IPAddresses = append(template.IPAddresses, net.ParseIP("127.0.0.1"))
	template.DNSNames = append(template.DNSNames, "localhost")

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(nil)

	pemBlock := &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}
	pem.Encode(buf, pemBlock)
	certPEM := append([]byte(nil), buf.Bytes()...)

	buf.Reset()

	pemBlock = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	pem.Encode(buf, pemBlock)
	keyPEM := append([]byte(nil), buf.Bytes()...)

	return &TLSConfig{CA: certPEM, Cert: certPEM, Key: keyPEM}
}
