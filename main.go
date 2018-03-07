package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	url     string
	method  string
	headers http.Header
	body    []byte
	debug   bool
}

type Desync struct {
	q     chan Message
	debug bool
}

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

	if m.debug {
		log.Printf("%T: %v\n", resp, resp)
	}

	return resp
}

func (d Desync) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.String() == "/" {
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer r.Body.Close()

	m := Message{r.URL.String()[1:], r.Method, r.Header, b, d.debug}

	if d.debug {
		log.Printf("%T: %v\n", m, m)
	}
	// send Message to channel
	d.q <- m
}

func (d Desync) serve(port string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Fatal(http.ListenAndServe(":"+port, d))
}

func (d *Desync) readChan(wg *sync.WaitGroup) {
	defer wg.Done()
	for m := range d.q {
		go m.send()
	}
}

func main() {

	var port string
	var wg sync.WaitGroup
	var d = Desync{make(chan Message), false}

	flag.BoolVar(&d.debug, "debug", false, "enable verbose logging")
	flag.StringVar(&port, "port", "8080", "port to listen")

	flag.Parse()

	wg.Add(2)
	go d.serve(port, &wg)
	go d.readChan(&wg)
	wg.Wait()
}
