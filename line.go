package httpoi

import (
	"io"
	"bytes"
)

var crlf = []byte("\r\n")

type LineWriter struct{
	w io.Writer
}

func NewLineWriter(w io.Writer) *LineWriter {
	return &LineWriter{
		w: w,
	}
}

func (lw *LineWriter) Write(data []byte) (int, error) {
	return lw.w.Write(append(data, '\r', '\n'))
}

func (lw *LineWriter) WriteEmpty() (int, error) {
	return lw.w.Write(crlf)
}

type LineReader struct{
	r io.Reader
	buf *bytes.Buffer
	b []byte
}

func NewLineReader(r io.Reader) *LineReader {
	return &LineReader{
		r: r,
		buf: bytes.NewBuffer(make([]byte, 0, 128)),
		b: make([]byte, 1, 1),
	}
}

func (lr *LineReader) Read() (content []byte) {
	for {
		lr.r.Read(lr.b)
		if lr.b[0] == '\r' {
			lr.r.Read(lr.b)
			if lr.b[0] == '\n' {
				content = lr.buf.Bytes()
				lr.buf.Reset()
				return
			}else{
				lr.buf.WriteByte('\r')
				lr.buf.Write(lr.b)
			}
		}else{
			lr.buf.Write(lr.b)
		}
	}
}