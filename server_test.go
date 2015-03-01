package httpoi

import (
	"testing"
	"strconv"
	"bytes"
	"fmt"
)

func Test_Sever(t *testing.T) {
	go Sever(":8181", func(c *Conn){
		c.Response.Headers["Content-Type"] = "text/html"
		html := []byte("Hello, World!")
		c.Response.Headers["Content-Length"] = strconv.Itoa(len(html))
		c.Response.StatusCode = StatusOK
		c.Response.Status()
		c.Response.Write(html)
	})

	Sever(":8182", func(c *Conn){
		c.Response.Headers["Content-Type"] = "text/html"
		html := bytes.NewBuffer([]byte("Hello, Chunked!"))
		c.Response.Headers["Transfer-Encoding"] = "chunked"
		c.Response.StatusCode = StatusOK
		c.Response.Status()
		ChunkedEncode(c.Response, html)
	})
}