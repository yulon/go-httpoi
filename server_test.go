package httpoi

import (
	"testing"
)

func Test_Sever(t *testing.T) {
	Sever(":8181", func(rs *ResponseW, rq *RequestR){
		rs.Headers["Content-Type"] = "text/html"
		rs.Headers["Transfer-Encoding"] = "chunked"
		rs.Headers["Content-Encoding"] = "gzip"
		rs.WriteHeaders(StatusOK)
		rs.WriteString("Hello, World!")
	})
}