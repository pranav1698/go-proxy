package proxy

import (
	"log"	
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	// "net/http/httputil"
	"io/ioutil"
	"sync"
	"github.com/pranav1698/go-proxy/config"
	"github.com/pranav1698/go-proxy/client"
)

type ProxyServer struct {
	Wg *sync.WaitGroup
	Port string
	RdbClient *client.RedisClient
}


func NewProxyServer(wg *sync.WaitGroup, conf config.Config) *ProxyServer {
	rdbClient, err := client.Connect(conf)
	if err != nil {
		log.Println(err)
		return nil
	}
	
	return &ProxyServer {
		Wg: wg,
		Port: conf.ProxyPort,
		RdbClient: rdbClient,
	}
}

func (ps ProxyServer) StartProxyServer() {
	defer ps.Wg.Done()
	log.Println("Starting Proxy Server")

	r := gin.Default()

	r.GET("/", ps.HandleProxy)

	r.Run(":9000")
}

func (ps ProxyServer) HandleProxy(c *gin.Context) {
	endpoint := c.Query("endpoint")
	if len(endpoint) == 0 {
		c.Data(http.StatusInternalServerError, "application/json", []byte("Endpoint not provided, please provide endpoint in query\n"))
		log.Println("Endpoint not provided")
		return
	}

	cacheData, err := ps.RdbClient.GetData(endpoint)
	if len(cacheData) != 0 {
		c.Data(http.StatusAccepted, "application/json", []byte(cacheData))
		log.Println("Date Entry found in Cache")
		return
	}

	log.Println("Request to endpoint: ", c.Request.Host)

	// creating a new http request
	req, err := http.NewRequest(c.Request.Method, endpoint, c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = c.Request.Header.Clone()
	req.URL.RawQuery = c.Request.URL.RawQuery

	// so that the request does not hang out
	client := http.Client{
		Timeout: 5*time.Second,
	}

	log.Println("Forwarding Request Data to ", endpoint)

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

	err = ps.RdbClient.AddData(endpoint, string(respBody))
	if err != nil {
		log.Println("Not able to enter data in redis cache")
	}

	c.Data(200, "application/json", respBody)
}