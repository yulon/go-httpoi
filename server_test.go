package httpoi

import (
	"testing"
)

func Test_Sever(t *testing.T) {
	Sever(":8181", func(r *Response){
		r.Headers["Content-Type"] = "text/html"
		r.WriteString("hello, world!")
	})
}