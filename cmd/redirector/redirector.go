package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jeremy-clerc/r"
)

var version string

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Version: %v\n", version)
	}
	var (
		path   = flag.String("path", "links", "Path to a file of shortcut!link.")
		listen = flag.String("listen", "127.0.0.1:8008", "Address and port to listen on.")
	)
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		log.Fatalf("failed to open links file: %v", err)
	}
	defer f.Close()

	rr, err := r.Load(f)
	if err != nil {
		log.Fatalf("failed to load content of links file: %v", err)
	}

	log.Fatal(http.ListenAndServe(*listen, rr))
}
