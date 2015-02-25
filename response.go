package httpoi

import (
	"net"
	"strconv"
	"runtime"
)

type Response struct{
	Version string
	Status string
	Headers map[string]string
	//////////////////////////
	conn net.Conn
	async bool
	close bool
}

var space = []byte(" ")
var crlf = []byte("\r\n")
var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")
var langVer = runtime.Version()

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