package mit

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

// ParsePrivateKey decodes the given private key DER block. OpenSSL 0.9.8
// generates PKCS#1 private keys by default, while OpenSSL 1.0.0 generates
// PKCS#8 keys.
func ParsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		return key, nil
	}
	return nil, errors.New("failed to parse private key")
}

type SubjectPublicKeyInfo struct {
	Algo      pkix.AlgorithmIdentifier
	BitString asn1.BitString
}

type PublicKeyAndChallenge struct {
	SubjectPublicKeyInfo SubjectPublicKeyInfo
	Challenge            string `asn1:"ia5"`
}

type SignedPublicKeyAndChallenge struct {
	PublicKeyAndChallenge PublicKeyAndChallenge
	SignatureAlgorithm    []asn1.ObjectIdentifier
	Signature             asn1.BitString
}

var oidSignatureMD5WithRSA = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 4}

func GetClientCertificate(username, mitid, password string, pk *rsa.PublicKey, sk *rsa.PrivateKey) (*x509.Certificate, error) {
	const ffx_useragent = "Mozilla/5.0 (X11; Linux x86_64; rv:23.0) Gecko/20100101 Firefox/23.0"
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// get a session cookie
	client := &http.Client{Jar: jar}
	req, err := http.NewRequest("GET", "https://ca.mit.edu/ca/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-agent", ffx_useragent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	// log in
	login := url.Values{}
	login.Set("data", "1")
	login.Set("login", username)
	login.Set("password", password)
	login.Set("mitid", mitid)
	login.Set("Submit", "Next >>")

	req, err = http.NewRequest("POST", "https://ca.mit.edu/ca/login", strings.NewReader(login.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-agent", ffx_useragent)
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	rr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ml := regexp.MustCompile("challenge=\"([0-9a-zA-Z]+)\"").FindSubmatch(rr)
	if len(ml) < 2 {
		return nil, fmt.Errorf("could not extract server challenge from the CA-s reply")
	}
	challenge := string(ml[1])

	var spki SubjectPublicKeyInfo
	pkBytes, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return nil, err
	}
	if _, err := asn1.Unmarshal(pkBytes, &spki); err != nil {
		return nil, err
	}

	var pkac PublicKeyAndChallenge
	pkac.SubjectPublicKeyInfo = spki
	pkac.Challenge = challenge

	pkacDer, err := asn1.Marshal(pkac)
	if err != nil {
		return nil, err
	}
	pkacHash := md5.Sum(pkacDer)

	var spkac SignedPublicKeyAndChallenge
	spkac.PublicKeyAndChallenge = pkac
	spkac.SignatureAlgorithm = append(spkac.SignatureAlgorithm, oidSignatureMD5WithRSA)
	spkac.Signature.Bytes, err = rsa.SignPKCS1v15(rand.Reader, sk, crypto.MD5, pkacHash[:])
	spkac.Signature.BitLength = 8 * len(spkac.Signature.Bytes)
	if err != nil {
		return nil, err
	}
	spkacDer, err := asn1.Marshal(spkac)
	if err != nil {
		return nil, err
	}

	// get the key certified
	post := url.Values{}
	post.Set("data", "1")
	post.Set("userkey", base64.StdEncoding.EncodeToString(spkacDer))
	post.Set("life", "1")
	post.Set("Submit", "Next >>")

	req, err = http.NewRequest("POST", "https://ca.mit.edu/ca/handlemoz", strings.NewReader(post.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-agent", ffx_useragent)
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	// retrieve the signed certificate
	req, err = http.NewRequest("GET", "https://ca.mit.edu/ca/mozcert/2", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-agent", ffx_useragent)
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	rr, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(rr)
}

/*
func main() {
	sk, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}
	keyPEMBlock, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return
	}

	cert, err := GetClientCertificate(os.Args[1], os.Args[2], os.Args[3], &sk.PublicKey, sk)
	if err != nil {
		log.Fatal(err)
	}

	pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
}
*/
