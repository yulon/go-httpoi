package httpoi

import (
	"io"
)

type RequestLine struct{
	Method string
	URI string
	HTTPVersion string
}

type RequestHeader struct{
	*RequestLine
	Fields HeaderFields
}

type RequestR struct{
	*RequestHeader
	io.Reader
}