package r

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestLoad(t *testing.T) {
	input := `m!https://mail.google.com/mail/u/0/?hl=fr#inbox
c!https://calendar.google.com/calendar/r
garbageline`
	entries := map[string]string{
		"m": "https://mail.google.com/mail/u/0/?hl=fr#inbox",
		"c": "https://calendar.google.com/calendar/r",
	}
	r, err := Load(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Load(%v): got err %v != want nil", input, err)
	}
	for shorcut, link := range entries {
		val, ok := r.Links.Load(shorcut)
		if !ok {
			t.Errorf("R.Links.Load(%v): not found", shorcut)
			continue
		}
		if val.(string) != link {
			t.Errorf("R.Links.Load(%v): got %v != want %v", shorcut, val.(string), link)
		}
	}
}

func TestServeHTTP(t *testing.T) {
	data := []struct {
		shorcut string
		loc     string
		code    int
	}{
		{"c", "https://calendar.google.com/calendar/r", http.StatusTemporaryRedirect},
		{"", "", http.StatusBadRequest},
		{"undef", "", http.StatusNotFound},
	}
	r := &R{Links: &sync.Map{}}
	r.Links.Store(data[0].shorcut, data[0].loc)

	srv := httptest.NewServer(r)
	defer srv.Close()

	cli := http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, d := range data {
		res, err := cli.Get(srv.URL + "/" + d.shorcut)
		if err != nil && err != http.ErrUseLastResponse {
			t.Errorf("GET /%v: got err %v != want nil/%v", d.shorcut, err, http.ErrUseLastResponse)
			continue
		}
		res.Body.Close()

		if res.StatusCode != d.code {
			t.Errorf("GET /%v: got status %v !=  want %v", d.shorcut, res.StatusCode, http.StatusTemporaryRedirect)
		}
		if loc := res.Header.Get("Location"); d.loc != "" && loc != d.loc {
			t.Errorf("GET /%v: got location %v != want %v", d.shorcut, loc, d.loc)
		}
	}
}
