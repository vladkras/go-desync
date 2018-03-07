package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {

	handler := &Desync{}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	log.Printf("%T: %v\n", handler.q, handler.q)
}

func TestMessage(t *testing.T) {
	m := &Message{"https://ya.ru", "GET", nil, nil, true}
	resp := m.send()
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
