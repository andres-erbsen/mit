package mit

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
	"os"
	"testing"
)

var username = os.Getenv("MIT_USER")
var mitid = os.Getenv("MIT_ID")
var password = os.Getenv("MIT_PASSWORD")

func TestGetClientCertificate(t *testing.T) {
	caCert, err := DownloadClientCA()
	if err != nil {
		t.Skip("failed to download client CA")
	}
	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)

	sk, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatal(err)
	}

	cert, err := GetClientCertificate(username, mitid, password, &sk.PublicKey, sk)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := cert.Verify(x509.VerifyOptions{Roots: caPool, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}}); err != nil {
		t.Fatal(err)
	}
}
