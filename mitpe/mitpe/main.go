package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/andres-erbsen/mit"
	"github.com/andres-erbsen/mit/mitpe"
)

var certfile = os.Getenv("MIT_CERT")
var skfile = os.Getenv("MIT_SK")
var mitid = os.Getenv("MIT_ID")

var course = flag.String("coursenumber", "", "course number. Example: 0608")
var day = flag.String("day", "", "when the section runs. Example: MW")
var time = flag.String("time", "", "when the section runs. Example (you probably want to quote it in the shell): 1:00 PM")

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	flag.Parse()
	if *course == "" || *day == "" || *time == "" || mitid == "" || certfile == "" || skfile == "" {
		fmt.Fprintf(os.Stderr, "USAGE: env MIT_CERT= MIT_SK= MIT_ID= %s -course=0608 -course=\"TR 2:00\"\n", os.Args[0])
		os.Exit(4)
	}
	cert, err := tls.LoadX509KeyPair(certfile, skfile)
	if err != nil {
		log.Fatal(err)
	}
	c := mit.NewClient(cert)
	err = mitpe.GetMePE(c, mitid, *course, *day, *time)
	if err != nil {
		log.Fatal(err)
	}
}
