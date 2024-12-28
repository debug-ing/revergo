package tls

import (
	"crypto/tls"
	"log"
)

// Load loads the SSL certificate and key from the provided files and returns a tls.Config.
func Load(certFile, keyFile string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load SSL certificate and key: %v", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}
