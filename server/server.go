package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

func main() {
	StartServer()
}

func StartServer() {
	r := gin.Default()

	r.GET("/", GetRequestHandler)

	r.Run(":9001")
}

func GetRequestHandler(c *gin.Context) {
	log.Println("Request to endpoint: ", c.Request.Host)

	c.Writer.Header().Add("test-header", "test-header-value")
	c.Writer.WriteHeader(http.StatusAccepted)
	c.Writer.Write([]byte("Server Endpoint Called\n"))
}