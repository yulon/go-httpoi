package httpoi

import (
	"bytes"
)

type Headers map[string]string

func (hf Headers) MakeHeaders() string {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	for k, v := range hf {
		buf.WriteString(k + ": "+ v + "\r\n")
	}
	return buf.String()
}