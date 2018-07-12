package stats_api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParameterPP01(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(Handler))
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Error by http.Get(). %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Error by ioutil.ReadAll(). %v", err)
	}

	newLine := strings.Count(string(data), "\n")
	if newLine != 0 {
		t.Fatalf("Data Error. %v / %d", string(data), newLine)
	}
}

func TestParameterPP02(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(Handler))
	defer ts.Close()

	for _, v := range []string{"1", "true"} {
		rp := fmt.Sprintf("?pp=%s", v)
		r, err := http.Get(ts.URL + rp)
		if err != nil {
			t.Fatalf("Error by http.Get(). %v", err)
		}
		data, _ := ioutil.ReadAll(r.Body)
		newLine := strings.Count(string(data), "\n")
		if newLine == 0 {
			t.Fatalf("pp parameter isn't work well.: %v/%d", string(data), newLine)
		}
	}
}
