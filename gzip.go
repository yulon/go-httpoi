package httpoi

import (
	"compress/gzip"
	"bytes"
	"io/ioutil"
)

func GzipCompress(data []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, len(data)))
	gz := gzip.NewWriter(buf)
	gz.Write(data)
	gz.Close()
	return buf.Bytes()
}

func GzipUncompress(data []byte) []byte {
	buf := bytes.NewBuffer(data)
	gz, _ := gzip.NewReader(buf)
	ret, _ := ioutil.ReadAll(gz)
	gz.Close()
	return ret
}