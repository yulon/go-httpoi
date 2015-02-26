package httpoi

type Request struct{
	Line RequestLine
	Headers map[string]string
	Body RequestBody
}

type RequestLine struct{
	Method string
	Uri RequestURI
	Version string
}

type RequestURI struct{
	Raw string
	Param map[string]string
}

type RequestBody struct{
	Raw []byte
	Param map[string]string
}