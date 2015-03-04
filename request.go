package httpoi

type RequestLine struct{
	Method string
	URI string
	HTTPVersion string
}

type RequestHeader struct{
	*RequestLine
	Fields map[string]string
}

type Requested struct{
	*RequestHeader
	Path string
	Params map[string]string
}