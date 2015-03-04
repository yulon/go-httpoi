package httpoi

import (
	"testing"
)

func Test_Sever(t *testing.T) {
	Sever(":8181", func(o *Respond, i *Requested){
		o.Fields["Content-Type"] = "text/html"
		o.Fields["Content-Encoding"] = "gzip"
		o.Status(StatusOK)
		o.WriteString("Hello, World!")
	})
}