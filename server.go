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
									req.Method = string(tmp)
								case 1:
									req.Url = string(tmp)
							}
							ix++
							tmp = nil
						}
					case '\r':
						req.Version = string(tmp)
						i++
						break m
					default:
						tmp = append(tmp, rawReq[i])
				}
			}

			if len(req.Version) > 4 && req.Version[:4] == "HTTP" { // Is HTTP
				tmp = nil
				var name string

				for ; i < rawReqLen; i++ { // Parse Request Headers
					switch rawReq[i] {
						case ':':
							if i+1 < rawReqLen && rawReq[i+1] == ' ' {
								name = string(tmp)
								tmp = nil
								i++
							}else{
								tmp = append(tmp, rawReq[i])
							}
						case '\r':
							if name != "" {
								req.Headers[name] = string(tmp)
								name = ""
							}
							tmp = nil
							i++
						default:
							tmp = append(tmp, rawReq[i])
					}
				}

				resp := &Response{
					Headers: map[string]string{
						"Date": time.Now().Format(time.RFC1123),
						"Sever": "HTTPOI",
						"X-Powered-By": langVer,
					},
					Version: req.Version,
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