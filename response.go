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

type ResponseHeader struct{
	*StatusLine
	Headers map[string]string
}

type ResponseW struct{
	*ResponseHeader
	wc io.WriteCloser
}

func (rs *ResponseW) WriteHeaders(StatusCode int) {
	rs.Status(StatusCode)
	lw := NewLineWriter(rs.wc)

	// Status Line
	lw.Write([]byte(rs.HTTPVersion + " " + strconv.Itoa(rs.StatusCode) + " " + rs.ReasonPhrase))

	// Headers
	for name, value := range rs.Headers {
		lw.Write([]byte(name + ": " + value)) // Header Field
	}
	lw.WriteEmpty() // Headers End

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