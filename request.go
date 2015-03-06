package httpoi

import (
	"io"
)

type RequestLine struct{
	Method string
	URI string
	HTTPVersion string
}

type RequestR struct{
	Line *RequestLine
	Fields HeaderFields
	io.Reader
}