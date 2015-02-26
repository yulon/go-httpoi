package httpoi

import (
	"net"
	"fmt"
	"bytes"
	"time"
)

func Sever(laddr string, handler func(*Response, Request)) {
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

func saw(c net.Conn, handler func(*Response, Request)) {
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

			req := Request{
				Headers: map[string]string{},
			}

			var i int
			var tmp []byte
			var ix int
			m:
			for i = 0; i < rawReqLen; i++ { // Parse Request Line
				switch rawReq[i] {
					case ' ':
						if tmp != nil {
							switch ix {
								case 0:
									req.Line.Method = string(tmp)
								case 1:
									req.Line.Uri.Raw = string(tmp)
							}
							ix++
							tmp = nil
						}
					case '\r':
						req.Line.Version = string(tmp)
						i++ // \n
						break m
					default:
						tmp = append(tmp, rawReq[i])
				}
			}

			if len(req.Line.Version) > 4 && req.Line.Version[:4] == "HTTP" { // Is HTTP
				tmp = nil
				var name string
				var colon bool

				for ; i < rawReqLen; i++ { // Parse Request Headers
					switch rawReq[i] {
						case ':':
							if !colon {
								colon = true
								name = string(tmp)
								tmp = nil
								for i+1 < rawReqLen && rawReq[i+1] == ' ' {
									i++
								}
							}

						case '\r':
							if name != "" {
								req.Headers[name] = string(tmp)
								name = ""
							}
							tmp = nil
							i++ // \n
							colon = false
							
						default:
							tmp = append(tmp, rawReq[i])
					}
				}

				resp := &Response{
					Headers: map[string]string{
						"Date": time.Now().Format(time.RFC1123),
						"Server": "HTTPOI",
						"X-Powered-By": langVer,
					},
					Version: req.Line.Version,
					conn: c,
				}

				handler(resp, req)

				if !resp.async {
					resp.Close()
				}
			}else{
				c.Close()
			}

			return
		}
	}
}