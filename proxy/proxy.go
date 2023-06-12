package proxy

import (
	"log"	
	"github.com/gin-gonic/gin"
	"github.com/pranav1698/go-proxy/config"
	"time"
	"net/http"
	// "net/http/httputil"
	"io/ioutil"
)

var conf config.Config = config.Config{
	ServerUrl: "http://localhost:9001",
}

func StartProxyServer() {
	log.Println("Starting Proxy Server")

	r := gin.Default()

	r.GET("/", HandleProxy)

	r.Run(":9000")
}

func HandleProxy(c *gin.Context) {
	log.Println("Request to endpoint: ", c.Request.Host)

	req, err := http.NewRequest(c.Request.Method, conf.ServerUrl, c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header = c.Request.Header.Clone()
	req.URL.RawQuery = c.Request.URL.RawQuery

	client := http.Client{
		Timeout: 5*time.Second,
	}

	log.Println("Forwarding Request Data to ", conf.ServerUrl)
	resp, err := client.Do(req)
	if err!= nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte("Error in forwarding request\n"))
		return
	}

	// respData, err := httputil.DumpResponse(resp, true)
	// if err != nil {
	// 	log.Println(err)
	// 	c.Writer.WriteHeader(http.StatusInternalServerError)
	// 	c.Writer.Write([]byte("Error in forwarding request\n"))
	// 	return
	// }

	// log.Println("Got Response from Server: ", string(respData))

	// for k, v := range resp.Header {
	// 	c.Writer.Header()[k] = v
	// }

	// _, err = io.Copy(c.Writer, resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	c.Writer.WriteHeader(http.StatusInternalServerError)
	// 	c.Writer.Write([]byte("Error in forwarding request\n"))
	// 	return
	// }

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "application/json", respBody)
}