package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jeremy-clerc/r/r"
)

func main() {
	var (
		path   = flag.String("path", "links", "Path to a file of shortcut!link.")
		listen = flag.String("listen", "127.0.0.1:8008", "Address and port to listen on.")
	)
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		log.Fatalf("failed to open links file: %v", err)
	}
	rr, err := r.Load(f)
	if err != nil {
		log.Fatalf("failed to load content of links file: %v", err)
	}

	log.Fatal(http.ListenAndServe(*listen, rr))
}
