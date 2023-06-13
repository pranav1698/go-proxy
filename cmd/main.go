package main

import (
	"github.com/pranav1698/go-proxy/proxy"
	"github.com/pranav1698/go-proxy/config"
	"sync"
)

func main() {
	var conf config.Config = config.Config{
		ServerUrl: "http://localhost:9001",
	}
	
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	
	ps := proxy.NewProxyServer(wg, conf)
	for i:=0; i<10; i++ {
		wg.Add(1)
		ps.StartProxyServer()
	}
	
	//ps.StartProxyServer()
}