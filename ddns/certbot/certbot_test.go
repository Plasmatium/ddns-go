package certbot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"testing"
)

func TestTLSExpiration(t *testing.T) {
	domain := "yucy-love.top"
	addr := domain + ":6443"

	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		t.Fatal("no certs found")
	}

	parsedCerts := make([]*x509.Certificate, 0, len(certs))
	for _, cert := range certs {
		parsed, err := x509.ParseCertificate(cert.Raw)
		if err != nil {
			panic(err)
		}

		parsedCerts = append(parsedCerts, parsed)
	}
	fmt.Println(parsedCerts)
}