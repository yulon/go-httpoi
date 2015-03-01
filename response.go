package httpoi

import (
	"runtime"
	"io"
)

type ResponseInfo struct{
	HTTPVersion string
	StatusCode string
	Headers map[string]string
}

type ResponseWriter struct{
	ResponseInfo
	wcr io.Writer
}

var space = []byte(" ")
var crlf = []byte("\r\n")
var langVer = runtime.Version()

func (resp ResponseWriter) writeHeader(content string) {
	resp.wcr.Write([]byte(content))
	resp.wcr.Write(crlf)
}

func (resp ResponseWriter) Status() {
	// line
	resp.wcr.Write([]byte(resp.HTTPVersion)) // Version version
	resp.wcr.Write(space)
	resp.wcr.Write([]byte(resp.StatusCode)) // status code
	resp.wcr.Write(crlf) // line end

	// headers
	for k, v := range resp.Headers {
		resp.writeHeader(k + ": "+ v)
	}
	resp.wcr.Write(crlf) // headers end
}

func (resp ResponseWriter) Write(data []byte) (int, error) {
	return resp.wcr.Write(data)
}