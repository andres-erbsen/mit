package mit

import (
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

func DownloadClientCA() (*x509.Certificate, error) {
	resp, err := http.Get("http://ca.mit.edu/mitClient.crt")
	if err != nil {
		return nil, err
	}
	rr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(rr)
}
