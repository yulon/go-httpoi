package httpoi

import (
	"testing"
)

func Test_Sever(t *testing.T) {
	Sever(":8181", func(resp *Response, req Request){
		resp.Headers["Content-Type"] = "text/html"
		resp.Write([]byte("hello, world!"))
	})
}