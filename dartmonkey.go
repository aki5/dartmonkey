package main

import (
	"crypto/rand"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var listenAddr = flag.String("listenAddr", ":8880", "address to listen on for http requests")
var dataFile = flag.String("dataFile", "sp500.csv", "file to read csv entries from (csv)")

type DartMonkey struct {
	records [][]string
}

func rand32() uint32 {
	var val uint32
	b := make([]byte, 4)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		log.Fatalf("rand32: %v\n", err)
	}
	val = uint32(b[0])
	val += uint32(b[1]) << 8
	val += uint32(b[2]) << 16
	val += uint32(b[3]) << 24
	return val
}

func (dart *DartMonkey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p2 := uint32(1)
	nrecords := uint32(len(dart.records) - 1)
	for p2 < nrecords {
		p2 *= 2
	}
	p2--
	rnd := rand32() & p2
	for rnd >= nrecords {
		rnd = rand32() & p2
	}
	record := dart.records[1+rnd]
	fmt.Fprintf(w, "Index: %v\n", 1+rnd)
	for i, rec := range record {
		fmt.Fprintf(w, "%v: %v\n", dart.records[0][i], rec)
	}
}

func main() {
	flag.Parse()
	file, err := os.Open(*dataFile)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Printf("%v\n", err)
		os.Exit(1)
	}
	file.Close()

	dartm := &DartMonkey{records: records}
	log.Printf("listening on %v\n", *listenAddr)
	err = http.ListenAndServe(*listenAddr, dartm)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
