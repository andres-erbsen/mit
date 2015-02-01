package mit

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func NewClient(cert tls.Certificate) *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	// tls.Certificate{Certificate: [][]byte{cert.Raw}, PrivateKey: sk, Leaf: cert}
	return &http.Client{
		Jar:       jar,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{cert}}}}
}

func TouchstoneLogin(c *http.Client, dst string) error {
	resp, err := c.Get(dst)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp, err = c.Get("https://idp.mit.edu:446/idp/Authn/Certificate?login_certificate=Use+Certificate+-+Go")
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	var SAMLResponse, RelayState string
	doc.Find("input[name=SAMLResponse]").Each(func(i int, s *goquery.Selection) {
		SAMLResponse, _ = s.Attr("value")
	})
	doc.Find("input[name=RelayState]").Each(func(i int, s *goquery.Selection) {
		RelayState, _ = s.Attr("value")
	})

	resp, err = c.PostForm("https://edu-apps.mit.edu/Shibboleth.sso/SAML2/POST", url.Values{"SAMLResponse": {SAMLResponse}, "RelayState": {RelayState}})
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return err
}
