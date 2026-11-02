package image

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	httppkg "github.com/thoas/picfit/http"
)

func TestFromURLRejectsDisallowedURLBeforeFetch(t *testing.T) {
	hit := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit = true
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	client := httppkg.NewClient(httppkg.WithFetchPolicy(httppkg.NewFetchPolicy(true, nil)))
	u, err := url.Parse(srv.URL + "/a.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := FromURL(context.Background(), u, client); !errors.Is(err, httppkg.ErrURLNotAllowed) {
		t.Fatalf("got %v", err)
	}
	if hit {
		t.Fatal("blocked literal was fetched")
	}

	client = httppkg.NewClient(httppkg.WithFetchPolicy(httppkg.NewFetchPolicy(false, []string{"cdn.example.com"})))
	u, err = url.Parse("http://evil.example/a.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := FromURL(context.Background(), u, client); !errors.Is(err, httppkg.ErrURLNotAllowed) {
		t.Fatalf("allow list got %v", err)
	}
}

func TestFromURLAllowsLoopbackWhenPolicyIsInactive(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = io.WriteString(w, "png")
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL + "/a.png")
	if err != nil {
		t.Fatal(err)
	}
	file, err := FromURL(context.Background(), u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	body, err := io.ReadAll(file.Stream)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "png" {
		t.Fatalf("body %q", body)
	}
}

func TestMimetypeDetectorsUseFetchPolicyClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fetched", http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL + "/a.png")
	if err != nil {
		t.Fatal(err)
	}
	client := httppkg.NewClient(httppkg.WithFetchPolicy(httppkg.NewFetchPolicy(true, nil)))

	for _, kind := range []string{MimetypeDetectorTypeSniff, MimetypeDetectorTypeFastimage} {
		detect := GetMimetypeDetectorWithClient(kind, client)
		if _, err := detect(u); !errors.Is(err, httppkg.ErrURLNotAllowed) {
			t.Fatalf("%s detector error = %v", kind, err)
		}
	}
}
