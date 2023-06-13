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
	wg := &sync.WaitGroup{}
	
	
	ps := proxy.NewProxyServer(wg, conf)
	for i:=0; i<10; i++ {
		ps.Wg.Add(1)
		ps.StartProxyServer()
	}
	
	ps.Wg.Wait()
}