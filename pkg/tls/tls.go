package tls

import (
	"crypto/tls"
	"log"
)

// TlsLoad loads the SSL certificate and key from the provided files
func TlsLoad(certFile, keyFile string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load SSL certificate and key: %v", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}
