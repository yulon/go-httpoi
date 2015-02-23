package httpoi

import (
	"net"
	"fmt"
	"bytes"
	"time"
	"strconv"
	"runtime"
)

var space []byte = []byte(" ")
var crlf []byte = []byte("\r\n")
var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")
var langVer = runtime.Version()
var serverVer = "HTTPOI"

var responseLineList map[int][]byte = map[int][]byte{
	200: []byte("200 OK"),
}

func Listen(laddr string) {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		fmt.Println("[HTTPOI:" + laddr + "] Failed to listen!")
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("[HTTPOI:" + laddr + "] Failed to get conn!")
		}
		go saw(conn)
	}
}

func saw(c net.Conn) {
	rawReqBuf := bytes.NewBuffer([]byte{})
	for { // reading Requestuest
		b := make([]byte, 512)
		leng, err := c.Read(b)
		if err != nil {
			return
		}
		rawReqBuf.Write(b[:leng])
		if leng <= 512 { // reading done
			rawReq := rawReqBuf.Bytes()
			rawReqLen := len(rawReq)
			fmt.Println(rawReq)

			req := &Request{
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

			if req.Version[:4] == "HTTP" {
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
								fmt.Println(name, string(tmp))
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
					ver: []byte(req.Version),
					conn: c,
					StatusCode: 200,
					Headers: map[string]string{
						"Sever": serverVer,
						"X-Powered-By": langVer,
						"Date": time.Now().Format(time.RFC1123),
						"Transfer-Encoding": "chunked",
					},
				}

				resp.ReturnHTML("hello, world!")

			}else{
				c.Close()
			}

			return
		}
	}
}

type Handler map[string]func(*Response, Request)

type route map[string]Handler

type Request struct{
	Method string
	Version string
	Url string
	Path string
	PathParam map[string]string
	GetParam map[string]string
	PostParam map[string]string
	Headers map[string]string
}

type Response struct{
	StatusCode int
	Headers map[string]string
	//////////////////////////
	conn net.Conn
	ver []byte
	ret bool
	async bool
}

func (resp Response) writeHeader(content string) {
	resp.conn.Write([]byte(content))
	resp.conn.Write(crlf)
}

func (resp Response) writeLahs() {
	// line
	resp.conn.Write(resp.ver) // HTTP version
	resp.conn.Write(space)
	resp.conn.Write(responseLineList[resp.StatusCode]) // status code
	resp.conn.Write(crlf) // line end

	// headers
	for k, v := range resp.Headers {
		resp.writeHeader(k + ": "+ v)
	}
	resp.conn.Write(crlf) // headers end
}

func (resp Response) write(content []byte) { // write chunk
	resp.conn.Write([]byte(strconv.FormatUint(uint64(len(content)), 16))) // size
	resp.conn.Write(crlf) // size end
	resp.conn.Write(content) // data
	resp.conn.Write(crlf) // data end
}

func (resp Response) close(){
	resp.conn.Write(lastChunkAndChunkedBodyEnd) // last-chunk + Chunked-Body end
	resp.conn.Close()
}

func (resp Response) ReturnHTML(html string) {
	resp.Headers["Content-Type"] = "text/html"
	resp.writeLahs()
	resp.write([]byte(html))
	resp.close()
}

func (resp Response) Async(){
	resp.async = true
}