package httpoi

import (
	"bytes"
)

type HeaderFields map[string]string

func (hf HeaderFields) MakeHF() string {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	for k, v := range hf {
		buf.WriteString(k + ": "+ v + "\r\n")
	}
	return buf.String()
}