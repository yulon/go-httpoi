package httpoi

import (
	"testing"
)

func Test_Sever(t *testing.T) {
	Sever(":8181", func(rs *ResponseW, rq *RequestR){
		rs.Status(StatusOK)
		rs.Fields["Content-Type"] = "text/html"
		rs.Fields["Transfer-Encoding"] = "chunked"
		rs.Fields["Content-Encoding"] = "gzip"
		rs.WriteHeader()
		rs.WriteString("Hello, World!")
	})
}