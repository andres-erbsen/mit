package mit

import (
	"crypto/tls"
	"os"
	"testing"
)

var certfile = os.Getenv("MIT_CERT")
var skfile = os.Getenv("MIT_SK")

func TestTouchstoneLogin(t *testing.T) {
	cert, err := tls.LoadX509KeyPair(certfile, skfile)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(cert)
	if err := TouchstoneLogin(c, "https://edu-apps.mit.edu/ose-rpt/?Search+Online+Reports=Search+Subject+Evaluation+Reports"); err != nil {
		t.Fatal(err)
	}
}
