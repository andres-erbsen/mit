package main

import (
	"crypto/rsa"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/andres-erbsen/mit"
)

func main() {
	username := os.Getenv("MIT_USER")
	mitid := os.Getenv("MIT_ID")
	password := os.Getenv("MIT_PASSWORD")
	if len(os.Args) != 2 || username == "" || mitid == "" || password == "" {
		fmt.Fprintf(os.Stderr, "USAGE: env MIT_USER= MIT_ID= MIT_PASSWORD= %s skfile\n", os.Args[0])
		os.Exit(2)
	}
	skPEMBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read secret key: %s\n", err)
		os.Exit(3)
	}

	skPEM, remainder := pem.Decode(skPEMBytes)
	if len(remainder) != 0 {
		fmt.Fprintf(os.Stderr, "secret key must be a single PEM block (%d bytes left over)", len(remainder))
		os.Exit(3)
	}

	sk, err := mit.ParsePrivateKey(skPEM.Bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse secret key: %s\n", err)
		os.Exit(3)
	}

	rsakey, ok := sk.(*rsa.PrivateKey)
	if !ok {
		fmt.Fprintf(os.Stderr, "only RSA keys are supported for MIT client certificates\n")
		os.Exit(3)
	}

	cert, err := mit.GetClientCertificate(username, mitid, password, &rsakey.PublicKey, rsakey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	if err = pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}); err != nil {
		fmt.Fprintf(os.Stderr, "writing to stdout: %s\n", err)
		os.Exit(4)
	}

}
