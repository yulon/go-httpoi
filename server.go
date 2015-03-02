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

			// Parse Request Line
			i, sp1, sp2 := 0, 0, 0
			prl:
			for ; i < rawReqLen; i++ {
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
						break prl
				}
			}

			if len(req.Line.HTTPVersion) > 4 && req.Line.HTTPVersion[:4] == "HTTP" { // Is HTTP
				req.Line.Method = string(rawReq[:sp1])
				req.Line.URI = string(rawReq[sp1+1:sp2])

				// Parse Request Headers
				req.Headers = map[string]string{}
				prh:
				for keyStart, keyEnd, valStart := i + 1, 0, 0; i < rawReqLen; i++ {
					switch rawReq[i] {
						case ':':
							if keyEnd == 0 {
								keyEnd = i
								for y := 1; i+1 < rawReqLen; y++ {
									if rawReq[i+y] != ' ' && rawReq[i+y] != '\t' {
										i = i + y - 1
										valStart = i + 1
										continue prh
									}
								}
							}

						case '\r':
							if keyEnd - keyStart > 0 {
								for y := 1; i+1 < rawReqLen; y++ {
									if rawReq[i-y] != ' ' && rawReq[i-y] != '\t' {
										req.Headers[string(rawReq[keyStart:keyEnd])] = string(rawReq[valStart:i-y+1])
										fmt.Println(string(rawReq[keyStart:keyEnd]), string(rawReq[valStart:i-y+1]))
										continue prh
									}
								}
							}

						case '\n':
							keyStart = i + 1
							keyEnd = 0
					}
				}

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