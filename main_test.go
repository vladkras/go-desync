package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestServer(t *testing.T) {

	handler := &Desync{}

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}

func TestDesyncServe(t *testing.T) {
	handler := &Desync{}

	wg := &sync.WaitGroup{}
	var port = "8080"
	var c = certs{}

	go handler.serve(port, c, wg)
	_, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", port), time.Second)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDesyncChan(t *testing.T) {

	d := &Desync{make(chan *Message, 1)}
	defer close(d.q)

	m := &Message{"https://example.com", "GET", nil, nil}

	d.q <- m

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go d.readChan(wg)

}

func TestMessage(t *testing.T) {
	m := &Message{"https://example.com", "GET", nil, nil}
	resp := m.send()
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
}
