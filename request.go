package httpoi

type RequestLine struct{
	Method string
	URI string
	HTTPVersion string
}

type RequestHeader struct{
	RequestLine
	Fields map[string]string
}

type RequestParser struct{
	RequestHeader
	Path string
	Params map[string]bos
}

type bos struct{
	data interface{}
}

func (bos bos) ToBytes() []byte {
	switch d := bos.data.(type){
		case []byte:
			return d
	}
	return []byte{}
}

func (bos bos) ToString() string {
	return string(bos.ToBytes())
}