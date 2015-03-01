package httpoi

type RequestInfo struct{
	Line RequestLine
	Headers map[string]string
}

type RequestLine struct{
	Method string
	URI string
	HTTPVersion string
}

type RequestParser struct{
	RequestInfo
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