package httpoi

import (
	"io"
	"errors"
	"net"
	//"fmt"
)

type RequestLine struct{
	Method string
	URI string
	HTTPVersion string
}

func ParseRequestLine(line []byte) (rl *RequestLine, err error) {
	rl = &RequestLine{}
	leng := len(line)
	for i := leng - 1; i >= 0; i-- {
		if line[i] == ' ' {
			rl.HTTPVersion = string(line[i+1:])
			if rl.HTTPVersion == "HTTP/1.1" {
				for y := 0; y < leng; y++ {
					if line[y] == ' ' {
						if y != i {
							rl.Method = string(line[:y])
							rl.URI = string(line[y+1:i])
						}else{
							err = errors.New("Request Line token count mismatch")
						}
						return
					}
				}
			}else{
				err = errors.New("Protocol is not supported")
				return
			}
		}
	}
	err = errors.New("Not Request")
	return
}

type RequestR struct{
	*RequestLine
	Headers map[string]string
	io.Reader
}

func ReadRequest(conn net.Conn) (*RequestR, error) {
	lr := NewLineReader(conn)

	rqLine, err := ParseRequestLine(lr.Read())
	if err != nil {
		return nil, err
	}

	headers := map[string]string{}
	for {
		line := lr.Read()
		if len(line) == 0 {
			return &RequestR{
				RequestLine: rqLine,
				Headers: headers,
				Reader: conn,
			}, nil
		}else{
			name, value := ParseHeaderField(line)
			if name != "" {
				headers[name] = value
			}
		}
	}
	return nil, errors.New("Header Field incomplete")
}