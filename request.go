package httpoi

type Request struct{
	Method string
	Version string
	Url string
	Path string
	PathParam map[string]string
	GetParam map[string]string
	PostParam map[string]string
	Headers map[string]string
}