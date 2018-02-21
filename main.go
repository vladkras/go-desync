package main

import (
	"flag"
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
var debug bool
var port string

func (m message) send() *http.Response {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(m.method, m.url, strings.NewReader(m.data.Encode()))
	if err != nil && debug {
		log.Printf("%v\n", err)
	}

	resp, err := client.Do(req)
	if err != nil && debug {
		log.Printf("%v\n", err)
	}

	return resp
}

func handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	params := mux.Vars(r)
	m := message{url: params["url"] + "?" + r.URL.Query().Encode(), data: r.Form, method: r.Method}

	if debug {
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
		resp := m.send()
		if debug {
			log.Printf("Response: %v\n", resp)
		}
	}
}
