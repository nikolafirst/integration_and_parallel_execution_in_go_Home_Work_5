package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

const proxyAddr string = "localhost:8082"

var (
	counter        int    = 0
	firstInstance  string = "localhost:8080"
	secondInstance string = "localhost:8081"
)

func main() {}

func handlerProxy(w http.ResponseWriter, r *http.Request) {

	textBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	text := string(textBytes)
	if counter == 0 {
		if _, err := http.Post(firstInstance, "text/plain", bytes.NewBuffer([]byte(text))); err != nil {
			log.Fatalln(err)
		}
		counter++
		return
	}

	if _, err := http.Post(secondInstance, "text/plain", bytes.NewBuffer([]byte(text))); err != nil {
		log.Fatalln(err)
	}
	counter--
}
