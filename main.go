package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/namsral/flag"
)

// Message for channel
type Message struct {
	url     string // /http(s)://url without 1st /
	method  string
	headers http.Header
	body    []byte
	debug   bool
}

// Desync main object
type Desync struct {
	q     chan *Message
	debug bool
}

// send Message body and headers to url using method
func (m *Message) send() *http.Response {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(m.method, m.url, bytes.NewReader(m.body))
	if err != nil {
		log.Printf("%v\n", err)
	}
	req.Header = m.headers

	// Prvent remote server from keeping connection alive
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%v\n", err)
		return resp
	}

	if m.debug {
		log.Printf("%T: %v\n", resp, resp)
	}

	// explicitly close body to avoid leaks
	resp.Body.Close()
	return resp
}

// ServeHTTP : Creates Message from incoming HTTP reuqest and pushes to channel
func (d Desync) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// nothing to send so return default response: 200 OK
	if r.URL.String() == "/" {
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%v\n", err)
	}
	defer r.Body.Close()

	m := &Message{r.URL.String()[1:], r.Method, r.Header, b, d.debug}

	if d.debug {
		log.Printf("%T: %v\n", m, m)
	}
	// send Message to channel
	d.q <- m
}

// HTTP Listner
func (d Desync) serve(port string, cert string, wg *sync.WaitGroup) {
	defer wg.Done()

	// check path to *.crt and *.key and use TLS if found
	if cert != "" {
		c := certs{path: cert}
		err := c.GetCerts()
		if err == nil {
			// start secured server
			log.Fatal(http.ListenAndServeTLS(":"+port, c.crt, c.key, d))
			return
		} else {
			log.Printf("%s", err)
		}
	}

	log.Fatal(http.ListenAndServe(":"+port, d))
}

// Recieves Message from channel and executes send()
func (d *Desync) readChan(wg *sync.WaitGroup) {
	defer wg.Done()
	for m := range d.q {
		go m.send()
	}
}

func main() {

	var port, cert string
	var wg sync.WaitGroup
	var d = Desync{make(chan *Message), false}

	flag.BoolVar(&d.debug, "debug", false, "enable verbose logging")
	flag.StringVar(&port, "port", "8080", "port to listen")
	flag.StringVar(&cert, "cert", "", "Path to .crt and .key")

	flag.Parse()

	// One for server and one for reader
	wg.Add(2)
	go d.serve(port, cert, &wg)
	go d.readChan(&wg)
	wg.Wait()
}
