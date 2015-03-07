package chunked

import (
	"io"
	"strconv"
)

var crlf = []byte("\r\n")
var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")

type Writer struct{
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (cw *Writer) Write(data []byte) (int, error) {
	return cw.w.Write(concat([]byte(strconv.FormatUint(uint64(len(data)), 16)), crlf, data, crlf))
}

func (cw *Writer) Close() (err error) {
	_, err = cw.w.Write(lastChunkAndChunkedBodyEnd)
	return
}

type GzipWriter struct{
	*Writer
	level int
}

func NewGzipWriter(w io.Writer) *GzipWriter {
	return &GzipWriter{
		Writer: NewWriter(w),
	}
}

func (czw *GzipWriter) Write(data []byte) (int, error) {
	return czw.Writer.Write(gz(data))
}