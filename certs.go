package main

import (
	"crypto/tls"
	"errors"
	"path/filepath"
)

type certs struct {
	key  string
	crt  string
	path string
}

func (c *certs) GetCerts() error {

	if c.path == "" {
		return errors.New("Path to certs not defined")
	}

	// check path for *.crt  and *.key
	crt, _ := filepath.Glob(c.path + "/*crt")
	key, _ := filepath.Glob(c.path + "/*key")

	if len(crt) == 0 || len(key) == 0 {
		return errors.New("Either cert or key file not found")
	}

	// check cert and key compatibility
	_, err := tls.LoadX509KeyPair(crt[0], key[0])
	if err != nil {
		return err
	}

	// assign cert files
	c.crt = crt[0]
	c.key = key[0]

	return nil
}
