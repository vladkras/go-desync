package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type message struct {
	data   url.Values
	url    string
	method string
}

var q = make(chan message)

func (m message) send() *http.Response {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(m.method, m.url, strings.NewReader(m.data.Encode()))
	if err != nil {
		log.Printf("%v\n", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%v\n", err)
	}

	return resp
}

func handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	params := mux.Vars(r)
	m := message{url: params["url"] + "?" + r.URL.Query().Encode(), data: r.Form, method: r.Method}

	log.Printf("Request: %v\n", m)
	q <- m
}

func serve() {
	r := mux.NewRouter().SkipClean(true)
	r.HandleFunc("/{url:.*}", handle)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	go serve()
	for m := range q {
		resp := m.send()
		log.Printf("Response: %v\n", resp)
	}
}
