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

func (sl *StatusLine) ToString() string {
	return sl.HTTPVersion + " " + strconv.Itoa(sl.StatusCode) + " " + sl.ReasonPhrase + "\r\n"
}

type ResponseW struct{
	Line *StatusLine
	Fields HeaderFields
	w io.Writer
	close bool
}

var crlf = []byte("\r\n")
var langVer = runtime.Version()
var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")

func (rsw *ResponseW) writeHeader() {
	rsw.Line.Status(StatusOK)
	rsw.w.Write([]byte(rsw.Line.ToString() + rsw.Fields.ToString() + "\r\n"))
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

func (rsw *ResponseW) Close() {
	if !rsw.close && rsw.Fields["Transfer-Encoding"] == "chunked" {
		rsw.close = true
		rsw.w.Write(lastChunkAndChunkedBodyEnd) // Last Chunk + Chunked Body End
	}
}