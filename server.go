package httpoi

import (
	"net"
	"fmt"
	"bytes"
	"time"
	"strconv"
	"runtime"
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

var space = []byte(" ")
var crlf = []byte("\r\n")
var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")
var langVer = runtime.Version()

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
					Status: StatusOK,
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
	Version string
	Status string
	Headers map[string]string
	//////////////////////////
	conn net.Conn
	async bool
	close bool
}

func (resp Response) writeHeader(content string) {
	resp.conn.Write([]byte(content))
	resp.conn.Write(crlf)
}

func (resp Response) writeHeaders() {
	// line
	resp.conn.Write([]byte(resp.Version)) // Version version
	resp.conn.Write(space)
	resp.conn.Write([]byte(resp.Status)) // status code
	resp.conn.Write(crlf) // line end

	// headers
	for k, v := range resp.Headers {
		resp.writeHeader(k + ": "+ v)
	}
	resp.conn.Write(crlf) // headers end
}

func (resp Response) Write(content []byte) { // write chunk
	if resp.Headers["Transfer-Encoding"] != "chunked" {
		resp.Headers["Transfer-Encoding"] = "chunked"
		resp.writeHeaders()
	}
	resp.conn.Write([]byte(strconv.FormatUint(uint64(len(content)), 16))) // size
	resp.conn.Write(crlf) // size end
	resp.conn.Write(content) // data
	resp.conn.Write(crlf) // data end
}

func (resp Response) Close(){
	if !resp.close {
		resp.close = true
		if resp.Headers["Transfer-Encoding"] == "chunked" {
			resp.conn.Write(lastChunkAndChunkedBodyEnd) // last-chunk + Chunked-Body end
		}
		resp.conn.Close()
	}
}

func (resp Response) Async(){
	resp.async = true
}