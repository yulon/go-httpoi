package httpoi

import (
	"net"
	"runtime"
)

type Response struct{
	Version string
	Headers map[string]string
	//////////////////////////
	conn net.Conn
	async bool
	close bool
}

var space = []byte(" ")
var crlf = []byte("\r\n")
var langVer = runtime.Version()

func (resp Response) writeHeader(content string) {
	resp.conn.Write([]byte(content))
	resp.conn.Write(crlf)
}

func (resp Response) Status(code string) {
	// line
	resp.conn.Write([]byte(resp.Version)) // Version version
	resp.conn.Write(space)
	resp.conn.Write([]byte(code)) // status code
	resp.conn.Write(crlf) // line end

	// headers
	for k, v := range resp.Headers {
		resp.writeHeader(k + ": "+ v)
	}
	resp.conn.Write(crlf) // headers end
}

func (resp Response) Write(data []byte) (int, error) {
	return resp.conn.Write(data)
}

func (resp Response) Close(){
	if !resp.close {
		resp.close = true
		resp.conn.Close()
	}
}

func (resp Response) Async(){
	resp.async = true
}