package test

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
	"math/big"
	"net"
	"time"
)

type TLSConfig struct {
	CA   []byte
	Cert []byte
	Key  []byte
}

func (this *TLSConfig) Certificate() (tlsCert tls.Certificate, err error) {
	tlsCert, err = tls.X509KeyPair(this.Cert, this.Key)
	if err != nil {
		return
	}

	x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return
	}
	tlsCert.Leaf = x509Cert
	return
}

var errNotPEMData = errors.New("not PEM data")

func (this *TLSConfig) CACertificate() (*x509.Certificate, error) {
	pemBlock, _ := pem.Decode(this.CA)
	if pemBlock == nil {
		return nil, errNotPEMData
	}
	return x509.ParseCertificate(pemBlock.Bytes)
}

func (this *TLSConfig) Client() *tls.Config {
	if this == nil {
		return nil
	}

	tlsCert, err := this.Certificate()
	if err != nil {
		panic(err)
	}

	caCert, err := this.CACertificate()
	if err != nil {
		panic(err)
	}
	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)

	// see http://tip.golang.org/doc/go1.3#major_library_changes
	// for more about ServerName and InsecureSkipVerify

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		RootCAs:            caPool,
		ClientCAs:          caPool,
		CipherSuites:       fastCipherSuites,
		NextProtos:         []string{"kafka"},
		InsecureSkipVerify: true,
	}

	return tlsConfig
}

func (this *TLSConfig) Server() *tls.Config {
	if this == nil {
		return nil
	}

	tlsCert, err := this.Certificate()
	if err != nil {
		panic(err)
	}

	caCert, err := this.CACertificate()
	if err != nil {
		panic(err)
	}
	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		RootCAs:      caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
		CipherSuites: fastCipherSuites,
		NextProtos:   []string{"kafka"},
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig
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

// Code from $GOROOT/src/pkg/crypto/tls/generate_cert.go

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
