package proxy

import (
	"log"	
	"github.com/gin-gonic/gin"
	"github.com/pranav1698/go-proxy/config"
	"time"
	"net/http"
	// "net/http/httputil"
	"io/ioutil"
	"sync"
)

type ProxyServer struct {
	wg *sync.WaitGroup
	conf config.Config
	port string
}


func NewProxyServer(wg *sync.WaitGroup, conf config.Config) ProxyServer {
	return ProxyServer {
		wg: wg,
		port: ":9000",
		conf: conf,
	}
}

func (ps ProxyServer) StartProxyServer() {
	defer ps.wg.Done()
	
	log.Println("Starting Proxy Server")

	r := gin.Default()

	r.GET("/", ps.HandleProxy)

	r.Run(":9000")
}

func (ps ProxyServer) HandleProxy(c *gin.Context) {
	log.Println("Request to endpoint: ", c.Request.Host)

	// creating a new http request
	req, err := http.NewRequest(c.Request.Method, ps.conf.ServerUrl, c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = c.Request.Header.Clone()
	req.URL.RawQuery = c.Request.URL.RawQuery

	// so that the request does not hang out
	client := http.Client{
		Timeout: 5*time.Second,
	}

	log.Println("Forwarding Request Data to ", ps.conf.ServerUrl)

	// making the http request
	resp, err := client.Do(req)
	if err!= nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte("Error in forwarding request\n"))
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "application/json", respBody)
}