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

func (sl *StatusLine) ToString() string {
	return sl.HTTPVersion + " " + strconv.Itoa(sl.StatusCode) + " " + sl.ReasonPhrase + "\r\n"
}

type ResponseW struct{
	Line *StatusLine
	Fields HeaderFields
	io.Writer
}

func (rs *ResponseW) WriteHeader() (err error) {
	_, err = rs.Write([]byte(rs.Line.ToString() + rs.Fields.ToString() + "\r\n"))
	return
}

func (rs *ResponseW) WriteString(data string) (int, error) {
	return rs.Write([]byte(data))
}