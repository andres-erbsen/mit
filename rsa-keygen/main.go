package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "USAGE: %s skfile\n", os.Args[0])
		os.Exit(2)
	}
	skfile := os.Args[1]

	if _, err := os.Stat(skfile); err == nil {
		fmt.Fprintf(os.Stderr, "%s already exists\n", skfile)
		os.Exit(1)
	}
	sk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(skfile, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(sk)}), 0600); err != nil {
		fmt.Fprintf(os.Stderr, "writing output: %s\n", err)
		os.Exit(1)
	}
}
