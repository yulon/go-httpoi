package httpoi

import (
	"runtime"
	"io"
	"strconv"
	"compress/gzip"
	"bytes"
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

type ResponseHeader struct{
	*StatusLine
	Fields map[string]string
}

type Respond struct{
	*ResponseHeader
	w io.Writer
	end bool
}

var crlf = []byte("\r\n")
var langVer = runtime.Version()
var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")

func (o *Respond) writeHeader() {
	buf := bytes.NewBuffer(make([]byte, 0, 512))

	buf.WriteString(o.HTTPVersion + " " + strconv.Itoa(o.StatusCode) + " " + o.ReasonPhrase + "\r\n") // Status Line

	// Header Fields
	for k, v := range o.Fields {
		buf.WriteString(k + ": "+ v + "\r\n") // Header Field
	}
	buf.Write(crlf) // Header End

	o.w.Write(buf.Bytes())
}

func (o *Respond) Write(data []byte) {
	if !o.end {
		if o.Fields["Transfer-Encoding"] != "chunked" {
			o.Fields["Transfer-Encoding"] = "chunked"
			o.writeHeader()
		}
		if o.Fields["Content-Encoding"] == "gzip" {
			buf := bytes.NewBuffer(make([]byte, 0, len(data)))
			gz := gzip.NewWriter(buf)
			gz.Write(data)
			gz.Close()
			data = buf.Bytes()
		}

		o.w.Write(concat([]byte(strconv.FormatUint(uint64(len(data)), 16)), crlf, data, crlf)) // Chunk
	}
}

func (o *Respond) WriteString(data string) {
	o.Write([]byte(data))
}

func (o *Respond) writeEnd() {
	if !o.end && o.Fields["Transfer-Encoding"] == "chunked" {
		o.end = true
		o.w.Write(lastChunkAndChunkedBodyEnd) // Last Chunk + Chunked Body End
	}
}