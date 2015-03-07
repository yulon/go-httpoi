package httpoi

import (
	"io"
	"strconv"
	"github.com/yulon/httpoi/chunked"
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
	Fields HeaderFields
}

func (rsh *ResponseHeader) MakeHeader() string {
	return rsh.StatusLine.MakeLine() + rsh.Fields.MakeFields() + "\r\n"
}

type ResponseW struct{
	*ResponseHeader
	wc io.WriteCloser
}

func (rs *ResponseW) WriteHeader(StatusCode int) (err error) {
	_, err = rs.wc.Write([]byte(rs.ResponseHeader.MakeHeader()))

	if rs.Fields["Transfer-Encoding"] == "chunked" {
		if rs.Fields["Content-Encoding"] == "gzip" {
			rs.wc = chunked.NewGzipWriter(rs.wc)
		}else{
			rs.wc = chunked.NewWriter(rs.wc)
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