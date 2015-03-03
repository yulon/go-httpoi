package httpoi

import (
	"testing"
)

func Test_Sever(t *testing.T) {
	Sever(":8181", func(c *Conn){
		c.Response.Fields["Content-Type"] = "text/html"
		c.Response.Fields["Content-Encoding"] = "gzip"
		c.Response.StatusCode = StatusOK
		c.Response.Write([]byte("Hello, World!"))
	})
}