package main

import (
	"net/http"
	"net/http/httputil"
	"log"
	"time"
)


func main() {
	client := http.Client{
		Timeout: 5*time.Second,
	}

	for i:=0; i<20; i++ {
		resp, err := client.Get("http://localhost:9000")
		if err != nil {
			log.Fatal(err)
		}

		respBody, err := httputil.DumpResponse(resp, true)
		if err!=nil {
			log.Fatal(err)
		}

		log.Println(string(respBody))
	}
}