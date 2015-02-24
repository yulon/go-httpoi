package httpoi

import (
	"net"
	"fmt"
	"bytes"
	"time"
	"strconv"
	"runtime"
)

func Sever(laddr string, handler func(*Response)) {
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
var serverVer = "HTTPOI"
var etagHat = TimeToUnixString(time.Now()) + "."

var responseLineList = map[int][]byte{
	200: []byte("200 OK"),
	304: []byte("304 OK"),
}

func saw(c net.Conn, handler func(*Response)) {
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
						req.Protocol = string(tmp)
						i++
						break m
					default:
						tmp = append(tmp, rawReq[i])
				}
			}

			if len(req.Protocol) > 4 && req.Protocol[:4] == "HTTP" { // Is HTTP
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

				r := &Response{
					Request: req,
					Headers: map[string]string{
						"Date": time.Now().Format(time.RFC1123),
						"Transfer-Encoding": "chunked",
						"Sever": serverVer,
						"X-Powered-By": langVer,
					},
					StatusCode: 200,
					data: bytes.NewBuffer([]byte{}),
					conn: c,
				}

				handler(r)

				if !r.async {
					r.Close()
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
	Protocol string
	Url string
	Path string
	PathParam map[string]string
	GetParam map[string]string
	PostParam map[string]string
	Headers map[string]string
}

type Response struct{
	Request Request
	StatusCode int
	Headers map[string]string
	//////////////////////////
	data *bytes.Buffer
	conn net.Conn
	async bool
	close bool
}

func (r Response) writeHeader(content string) {
	r.conn.Write([]byte(content))
	r.conn.Write(crlf)
}

func (r Response) writeHeaders() {
	// line
	r.conn.Write([]byte(r.Request.Protocol)) // Protocol version
	r.conn.Write(space)
	r.conn.Write(responseLineList[r.StatusCode]) // status code
	r.conn.Write(crlf) // line end

	// headers
	for k, v := range r.Headers {
		r.writeHeader(k + ": "+ v)
	}
	r.conn.Write(crlf) // headers end
}

func (r Response) writeData() {
	if r.data.Len() != 0 {

		eTag := etagHat + "-" + strconv.Itoa(hash(r.data.Bytes()))

		if r.Request.Headers["If-None-Match"] == eTag {
			r.StatusCode = 304
			r.writeHeaders()
		}else{
			r.Headers["Etag"] = eTag
			r.Headers["Transfer-Encoding"] = "chunked"
			r.writeHeaders()
			r.writeChunkedBody(r.data)
		}

	}else{
		r.Headers["Content-Length"] = "0"
		r.writeHeaders()
	}
}

var chunkSize = []byte(strconv.FormatUint(uint64(5120), 16))

func (r Response) writeChunkedBody(buf *bytes.Buffer) {
	for { // write chunk
		data := make([]byte, 5120)
		leng, err := buf.Read(data)
		
		if err != nil {
			break
		}

		if leng < 5120{
			r.conn.Write([]byte(strconv.FormatUint(uint64(leng), 16))) // size
			r.conn.Write(crlf) // size end
			r.conn.Write(data[:leng]) // data
			r.conn.Write(crlf) // data end
			break
		}else{
			r.conn.Write(chunkSize) // size
			r.conn.Write(crlf) // size end
			r.conn.Write(data) // data
			r.conn.Write(crlf) // data end
		}
	}

	r.conn.Write(lastChunkAndChunkedBodyEnd) // last-chunk + Chunked-Body end
}

func (r Response) Write(content []byte) {
	r.data.Write(content)
}

func (r Response) WriteString(content string) {
	r.data.WriteString(content)
}

func (r Response) Close(){
	if !r.close {
		r.close = true
		r.writeData()
		r.conn.Close()
	}
}

func (r Response) Async(){
	r.async = true
}