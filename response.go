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

func (sl *StatusLine) Status(code int) {
	sl.StatusCode = code
	sl.ReasonPhrase = ReasonPhrases[sl.StatusCode]
}

type ResponseHeader struct{
	*StatusLine
	Fields map[string]string
}

type ResponseW struct{
	*ResponseHeader
	w io.Writer
	close bool
}

var crlf = []byte("\r\n")
var langVer = runtime.Version()
var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")

func (rsw *ResponseW) writeHeader() {
	buf := bytes.NewBuffer(make([]byte, 0, 512))

	buf.WriteString(rsw.HTTPVersion + " " + strconv.Itoa(rsw.StatusCode) + " " + rsw.ReasonPhrase + "\r\n") // Status Line

	// Header Fields
	for k, v := range rsw.Fields {
		buf.WriteString(k + ": "+ v + "\r\n") // Header Field
	}
	buf.Write(crlf) // Header End

	rsw.w.Write(buf.Bytes())
}

func (rsw *ResponseW) Write(data []byte) {
	if !rsw.close {
		if rsw.Fields["Transfer-Encoding"] != "chunked" {
			rsw.Fields["Transfer-Encoding"] = "chunked"
			rsw.writeHeader()
		}
		if rsw.Fields["Content-Encoding"] == "gzip" {
			buf := bytes.NewBuffer(make([]byte, 0, len(data)))
			gz := gzip.NewWriter(buf)
			gz.Write(data)
			gz.Close()
			data = buf.Bytes()
		}

		rsw.w.Write(concat([]byte(strconv.FormatUint(uint64(len(data)), 16)), crlf, data, crlf)) // Chunk
	}
}

func (rsw *ResponseW) WriteString(data string) {
	rsw.Write([]byte(data))
}

func (rsw *ResponseW) writeEnd() {
	if !rsw.close && rsw.Fields["Transfer-Encoding"] == "chunked" {
		rsw.close = true
		rsw.w.Write(lastChunkAndChunkedBodyEnd) // Last Chunk + Chunked Body End
	}
}