package httpoi

import (
	"io"
	"strconv"
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

type ResponseW struct{
	line *StatusLine
	Headers Headers
	wc io.WriteCloser
}

func (rs *ResponseW) Start(StatusCode int) {
	lw := NewLineWriter(rs.wc)

	// Status Line
	rs.line.Status(StatusCode)
	lw.Write([]byte(rs.line.HTTPVersion + " " + strconv.Itoa(rs.line.StatusCode) + " " + rs.line.ReasonPhrase))

	// Headers
	rs.Headers.WriteTo(lw)

	if rs.Headers["Transfer-Encoding"] == "chunked" {
		if rs.Headers["Content-Encoding"] == "gzip" {
			rs.wc = NewChunkedGzipWriter(rs.wc)
		}else{
			rs.wc = NewChunkedWriter(rs.wc)
		}
	}
}

func (rs *ResponseW) Write(data []byte) (int, error) {
	return rs.wc.Write(data)
}

func (rs *ResponseW) WriteString(data string) (int, error) {
	return rs.wc.Write([]byte(data))
}