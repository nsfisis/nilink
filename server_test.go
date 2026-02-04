package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(t *testing.T) *httptest.Server {
	t.Helper()
	db, err := openDB(":memory:")
	if err != nil {
		t.Fatalf("openDB: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	insertLink(db, "https://example.com")
	return httptest.NewServer(newMux(db))
}

func TestRobotsTxt(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/robots.txt")
	if err != nil {
		t.Fatalf("GET /robots.txt: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	ct := resp.Header.Get("Content-Type")
	if ct != "text/plain" {
		t.Errorf("Content-Type = %q, want text/plain", ct)
	}
	body, _ := io.ReadAll(resp.Body)
	want := "User-agent: *\nDisallow: /\n"
	if string(body) != want {
		t.Errorf("body = %q, want %q", body, want)
	}
}

func TestRedirect(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	short, _ := encodeID(1)
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Get(srv.URL + "/" + short)
	if err != nil {
		t.Fatalf("GET /%s: %v", short, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 301 {
		t.Errorf("status = %d, want 301", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	if loc != "https://example.com" {
		t.Errorf("Location = %q, want %q", loc, "https://example.com")
	}
}

func TestNotFoundInvalidID(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/ZZZZ")
	if err != nil {
		t.Fatalf("GET /ZZZZ: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestNotFoundRoot(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestNotFoundNested(t *testing.T) {
	srv := testServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/nested/path")
	if err != nil {
		t.Fatalf("GET /nested/path: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}
