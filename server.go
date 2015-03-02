package httpoi

import (
	"net"
	"fmt"
	"bytes"
	"time"
)

func Sever(laddr string, handler func(*Conn)) {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		fmt.Println("[HTTPOI:" + laddr + "] Failed to listen!")
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("[HTTPOI:" + laddr + "] Failed to get conn!")
		}
		go saw(conn, handler)
	}
}

type Conn struct{
	Request RequestParser
	Response *ResponseWriter
}

func saw(c net.Conn, handler func(*Conn)) {
	rawReqBuf := bytes.NewBuffer([]byte{})
	for { // reading Request
		b := make([]byte, 512)
		leng, err := c.Read(b)
		if err != nil {
			return
		}
		rawReqBuf.Write(b[:leng])
		if leng <= 512 { // reading done
			rawReq := rawReqBuf.Bytes()
			rawReqLen := len(rawReq)
			req := RequestParser{}
			i, sp1, sp2 := 0, 0, 0

			m:
			for ; i < rawReqLen; i++ { // Parse Request Line
				switch rawReq[i] {
					case ' ':
						if sp1 == 0 {
							sp1 = i
						}else{
							sp2 = i
						}

					case '\r':
						req.Line.HTTPVersion = string(rawReq[sp2+1:i])

					case '\n':
						break m
				}
			}

			if len(req.Line.HTTPVersion) > 4 && req.Line.HTTPVersion[:4] == "HTTP" { // Is HTTP
				req.Line.Method = string(rawReq[:sp1])
				req.Line.URI = string(rawReq[sp1+1:sp2])
				fmt.Println(req.Line)
				/*req.Headers = map[string]string{}

				start = 0
				var name string
				var colon bool

				for ; i < rawReqLen; i++ { // Parse Request Headers
					switch rawReq[i] {
						case ':':
							if !colon {
								colon = true
								name = string(rawReq[start:i])
								for i+1 < rawReqLen && rawReq[i+1] == ' ' {
									i++
								}
								start = i + 1
							}

						case '\r':
							if name != "" {
								req.Headers[name] = string(rawReq[start:i])
								name = ""
							}
							colon = false

						case '\n':
							start = i + 1
					}
				}*/

				resp := &ResponseWriter{
					ResponseInfo: ResponseInfo{
						Headers: map[string]string{
							"Date": time.Now().Format(time.RFC1123),
							"Server": "HTTPOI",
							"X-Powered-By": langVer,
						},
						HTTPVersion: req.Line.HTTPVersion,
					},
					wcr: c,
				}

				handler(&Conn{req, resp})
			}
			
			c.Close()
			return
		}
	}
}