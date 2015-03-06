package httpoi

import (
	"testing"
)

func Test_Sever(t *testing.T) {
	Sever(":8181", func(w *ResponseW, r *RequestR){
		w.Fields["Content-Type"] = "text/html"
		w.Fields["Content-Encoding"] = "gzip"
		w.WriteString("Hello, World!")
	})
}