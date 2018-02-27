package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Message struct {
	url     string
	method  string
	headers http.Header
	body    []byte
}

type Dummy struct{}

var q = make(chan Message)
var debug bool

func (m *Message) send() *http.Response {
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

func (d Dummy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.String() == "/" {
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer r.Body.Close()

	m := Message{r.URL.String()[1:], r.Method, r.Header, b}

	if debug {
		log.Printf("Request: %v\n", m)
	}
	// send Message to channel
	q <- m
}

func serve(port string) {
	log.Fatal(http.ListenAndServe(":"+string(port), Dummy{}))
	log.Println("Start...")
}

func main() {

	var port string

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
