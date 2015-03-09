package httpoi

import (
	"io"
	"strconv"
)

var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")

type ChunkedWriter struct{
	w io.Writer
}

func NewChunkedWriter(w io.Writer) *ChunkedWriter {
	return &ChunkedWriter{
		w: w,
	}
}

func (cw *ChunkedWriter) Write(data []byte) (int, error) {
	return cw.w.Write(concat([]byte(strconv.FormatUint(uint64(len(data)), 16)), crlf, data, crlf))
}

func (cw *ChunkedWriter) Close() (err error) {
	_, err = cw.w.Write(lastChunkAndChunkedBodyEnd)
	return
}

type ChunkedGzipWriter struct{
	*ChunkedWriter
	level int
}

func NewChunkedGzipWriter(w io.Writer) *ChunkedGzipWriter {
	return &ChunkedGzipWriter{
		ChunkedWriter: NewChunkedWriter(w),
	}
}

func (czw *ChunkedGzipWriter) Write(data []byte) (int, error) {
	return czw.ChunkedWriter.Write(gz(data))
}