package httpoi

import (
	"errors"
)

func ParseHeaderField(line string) (string, string) {
	for i := 0; i < len(line); i++ {
		if line[i] == ':' {
			return line[:i], trim(line[i+1:])
		}
	}
	return "", ""
}

type Headers map[string]string

func ParseHeaders(lr *LineReader) (hs Headers, err error) {
	hs = map[string]string{}
	for {
		line := lr.Read()
		if len(line) == 0 {
			return
		}else{
			name, value := ParseHeaderField(string(line))
			if name != "" {
				hs[name] = value
			}else{
				err = errors.New("Header Field incomplete")
				return
			}
		}
	}
}

func (hs Headers) WriteTo(lw *LineWriter) {
	for name, value := range hs {
		lw.Write([]byte(name + ": " + value)) // Header Field
	}
	lw.WriteEmpty() // Headers End
}