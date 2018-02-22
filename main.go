package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type message struct {
	body    []byte
	url     string
	method  string
	headers http.Header
}

var q = make(chan message)
var debug bool
var port string

func (m message) send() *http.Response {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(m.method, m.url, bytes.NewReader(m.body))
	if err != nil {
		log.Printf("%v\n", err)
	}
	req.Header = m.headers

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%v\n", err)
	}

	return resp
}

func (m message) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer r.Body.Close()

	m.url = r.URL.String()[1:]
	m.body = b
	m.method = r.Method
	m.headers = r.Header

	if debug {
		log.Printf("Request: %v\n", m)
	}
	// send message to channel
	q <- m
}

func serve(port string) {
	m := message{}
	log.Fatal(http.ListenAndServe(":"+string(port), m))
}

func main() {

	flag.BoolVar(&debug, "debug", false, "enable verbose logging")
	flag.StringVar(&port, "port", "8080", "port to listen")

	flag.Parse()

	go serve(port)
	for m := range q {
		r := m.send()
		if debug {
			log.Printf("Response: %v\n", r)
		}
	}
}
