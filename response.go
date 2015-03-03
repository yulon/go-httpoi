package httpoi

import (
	"runtime"
	"io"
	"strconv"
	"compress/gzip"
	"bytes"
)

type StatusLine struct{
	HTTPVersion string
	StatusCode int
	ReasonPhrase string
}

func (sl *StatusLine) StatusText() string {
	if sl.ReasonPhrase == "" {
		sl.ReasonPhrase = ReasonPhrases[sl.StatusCode]
	}
	return strconv.Itoa(sl.StatusCode) + sl.ReasonPhrase
}

type ResponseHeader struct{
	StatusLine
	Fields map[string]string
}

type ResponseWriter struct{
	ResponseHeader
	w io.Writer
}

var space = []byte(" ")
var crlf = []byte("\r\n")
var langVer = runtime.Version()

func (this ResponseWriter) writeField(content string) {
	this.w.Write([]byte(content))
	this.w.Write(crlf)
}

func (this ResponseWriter) writeHeader() {
	// Line
	this.w.Write([]byte(this.HTTPVersion))
	this.w.Write(space)
	this.w.Write([]byte(this.StatusText()))
	this.w.Write(crlf) // Line End

	// Fields
	for k, v := range this.Fields {
		this.writeField(k + ": "+ v)
	}
	this.w.Write(crlf) // Fields End
}

func (this ResponseWriter) Write(data []byte) {
	if this.Fields["Transfer-Encoding"] != "chunked" {
		this.Fields["Transfer-Encoding"] = "chunked"
		this.writeHeader()
	}
	if this.Fields["Content-Encoding"] == "gzip" {
		buf := bytes.NewBuffer([]byte{})
		gz := gzip.NewWriter(buf)
		gz.Write(data)
		gz.Close()
		data = buf.Bytes()
	}
	// chunk
	this.w.Write([]byte(strconv.FormatUint(uint64(len(data)), 16))) // size
	this.w.Write(crlf) // size end
	this.w.Write(data) // data
	this.w.Write(crlf) // data end
}