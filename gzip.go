package httpoi

import (
	"compress/gzip"
	"bytes"
	"io/ioutil"
)

func gz(data []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, len(data)))
	z := gzip.NewWriter(buf)
	z.Write(data)
	z.Close()
	return buf.Bytes()
}

func ungz(data []byte) []byte {
	buf := bytes.NewBuffer(data)
	z, _ := gzip.NewReader(buf)
	ret, _ := ioutil.ReadAll(z)
	z.Close()
	return ret
}