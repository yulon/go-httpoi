package httpoi

import (
	"bytes"
)

type HeaderFields map[string]string

func (hfs HeaderFields) MakeFields() string {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	for k, v := range hfs {
		buf.WriteString(k + ": "+ v + "\r\n")
	}
	return buf.String()
}