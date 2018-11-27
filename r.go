// Package r provides helpers to redirect users from shorcuts to defined URLs.
package r

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const delimiter byte = '!'

// R redirects shorcuts to defined URLs.
type R struct {
	Links *sync.Map
}

// ServeHTTP redirect users if we know about their shortcut.
func (r *R) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if len(req.URL.Path) < 2 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	val, ok := r.Links.Load(req.URL.Path[1:])
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, req, val.(string), http.StatusTemporaryRedirect)
}

// Load parses a text formated input and initialize a R handler.
// Expected lines format: shortcut!url
func Load(input io.Reader) (*R, error) {
	r := &R{Links: &sync.Map{}}
	s := bufio.NewScanner(input)
	for s.Scan() {
		b := s.Bytes()
		i := bytes.IndexByte(b, delimiter)
		if i == -1 || i-1 == len(b) {
			log.Printf("invalid line %v", s.Text())
			continue
		}
		r.Links.Store(string(b[:i]), string(b[i+1:]))
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("failed reading input: %v", err)
	}
	return r, nil
}
