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

func (sl *StatusLine) MakeLine() string {
	return sl.HTTPVersion + " " + strconv.Itoa(sl.StatusCode) + " " + sl.ReasonPhrase + "\r\n"
}

type ResponseHeader struct{
	*StatusLine
	Headers Headers
}

func (rsh *ResponseHeader) MakeHeaders() string {
	return rsh.StatusLine.MakeLine() + rsh.Headers.MakeHeaders() + "\r\n"
}

type ResponseW struct{
	*ResponseHeader
	wc io.WriteCloser
}

func (rs *ResponseW) WriteHeaders(StatusCode int) (err error) {
	rs.ResponseHeader.StatusLine.Status(StatusCode)
	_, err = rs.wc.Write([]byte(rs.ResponseHeader.MakeHeaders()))

	if rs.Headers["Transfer-Encoding"] == "chunked" {
		if rs.Headers["Content-Encoding"] == "gzip" {
			rs.wc = NewChunkedGzipWriter(rs.wc)
		}else{
			rs.wc = NewChunkedWriter(rs.wc)
		}
	}

	return
}

func (rs *ResponseW) Write(data []byte) (int, error) {
	return rs.wc.Write(data)
}

func (rs *ResponseW) WriteString(data string) (int, error) {
	return rs.wc.Write([]byte(data))
}