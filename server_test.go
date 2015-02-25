package httpoi

import (
	"testing"
	"strconv"
	"bytes"
)

func Test_Sever(t *testing.T) {
	go Sever(":8181", func(resp *Response, req Request){
		resp.Headers["Content-Type"] = "text/html"
		html := []byte("Hello, World!")
		resp.Headers["Content-Length"] = strconv.Itoa(len(html))
		resp.Status(StatusOK)
		resp.Write(html)
	})

	Sever(":8182", func(resp *Response, req Request){
		resp.Headers["Content-Type"] = "text/html"
		html := bytes.NewBuffer([]byte("Hello, Chunked!"))
		resp.Headers["Transfer-Encoding"] = "chunked"
		resp.Status(StatusOK)
		ChunkedEncode(resp, html)
	})
}