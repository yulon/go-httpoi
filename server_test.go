package httpoi

import (
	"testing"
)

func Test_Sever(t *testing.T) {
	Sever(":8181", func(rs *ResponseW, rq *RequestR){
		rs.Line.Status(StatusFound)
		rs.Fields["Location"] = "http://www.baidu.com"
		rs.WriteHeader()
	})
}