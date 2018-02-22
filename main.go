package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func handle(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	u := params["url"]

	if len(r.URL.Query()) > 0 {
		u = u + "?" + r.URL.Query().Encode()
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer r.Body.Close()

	m := message{url: u, body: b, method: r.Method, headers: r.Header}

	if debug {
		log.Printf("Recived Body: %s\n", b)
		log.Printf("Request: %v\n", m)
	}
	// send message to channel
	q <- m
}

func serve(port string) {
	r := mux.NewRouter().SkipClean(true)
	r.HandleFunc("/{url:.*}", handle)
	log.Fatal(http.ListenAndServe(":"+string(port), r))
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
			b, _ := ioutil.ReadAll(r.Body)
			log.Printf("Body: %s\n", b)
		}
	}
}
