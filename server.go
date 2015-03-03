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
	for { // Reading Request
		b := make([]byte, 512)
		leng, err := c.Read(b)
		if err != nil {
			return
		}
		rawReqBuf.Write(b[:leng])
		if leng <= 512 { // Reading Done
			rawReq := rawReqBuf.Bytes()
			fmt.Println(rawReq)
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
						req.HTTPVersion = string(rawReq[sp2+1:i])

					case '\n':
						break prl
				}
			}

			if len(req.HTTPVersion) > 4 && req.HTTPVersion[:4] == "HTTP" { // Is HTTP
				req.Method = string(rawReq[:sp1])
				req.URI = string(rawReq[sp1+1:sp2])

				// Parse Request Header Fields
				req.Fields = map[string]string{}
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
										req.Fields[string(rawReq[keyStart:keyEnd])] = string(rawReq[valStart:i-y+1])
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
					ResponseHeader: ResponseHeader{
						Fields: map[string]string{
							"Date": time.Now().Format(time.RFC1123),
							"Server": "HTTPOI",
							"X-Powered-By": langVer,
						},
					},
					w: c,
				}
				resp.HTTPVersion = req.HTTPVersion

				handler(&Conn{req, resp})

				if resp.Fields["Transfer-Encoding"] == "chunked" {
					resp.w.Write(lastChunkAndChunkedBodyEnd) // last-chunk + Chunked-Body end
				}
			}
			
			c.Close()
			return
		}
	}
}