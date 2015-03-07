package httpoi

import (
	"testing"
)

func Test_Sever(t *testing.T) {
	Sever(":8181", func(rs *ResponseW, rq *RequestR){
		rs.Status(StatusOK)
		rs.HF["Content-Type"] = "text/html"
		rs.HF["Transfer-Encoding"] = "chunked"
		rs.HF["Content-Encoding"] = "gzip"
		rs.WriteHeader()
		rs.WriteString("Hello, World!")
	})
}