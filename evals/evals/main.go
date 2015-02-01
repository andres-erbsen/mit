package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/andres-erbsen/mit"
	"github.com/andres-erbsen/mit/evals"
)

var certfile = os.Getenv("MIT_CERT")
var skfile = os.Getenv("MIT_SK")

var subject = flag.String("subject", "", "subject id")
var instructor = flag.String("instructor", "", "instructor name")
var department = flag.String("department", "", "department id")
var term = flag.String("term", "", "term id")

func main() {
	flag.Parse()
	if *subject == "" && *instructor == "" && *department == "" && *term == "" || certfile == "" || skfile == "" {
		fmt.Fprintf(os.Stderr, "USAGE: env MIT_CERT= MIT_SK= %s -subject=6.01\n", os.Args[0])
	}

	cert, err := tls.LoadX509KeyPair(certfile, skfile)
	if err != nil {
		log.Fatal(err)
	}
	c := mit.NewClient(cert)
	if err := mit.TouchstoneLogin(c, "https://edu-apps.mit.edu/ose-rpt/subjectEvaluationSearch.htm"); err != nil {
		log.Fatal(err)
	}
	m, err := evals.Search(c, *term, *department, *subject, *instructor)
	if err != nil {
		log.Fatal(err)
	}

	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(ks)))

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, ' ', 0)
	for _, k := range ks {
		mm, err := evals.Report(c, m[k])
		if err != nil {
			log.Fatal(err)
		}
		for kk, vv := range mm {
			fmt.Fprintf(w, "%s\t | %s\t | %s\n", k, kk, vv)
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
}
