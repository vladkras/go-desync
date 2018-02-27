package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessage(t *testing.T) {
	handler := &Dummy{}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
