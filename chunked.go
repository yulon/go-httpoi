package httpoi

import (
	"io"
	"strconv"
)

var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")

func Chunked(dest io.Writer, src io.Reader, size int) {
	var chunkSize = []byte(strconv.FormatUint(uint64(size), 16))

	for {
		data := make([]byte, size)
		leng, err := src.Read(data)

		if err != nil {
			break
		}

		if leng < size{
			// (last-1)-chunk
			dest.Write([]byte(strconv.FormatUint(uint64(leng), 16))) // size
			dest.Write(crlf) // size end
			dest.Write(data[:leng]) // data
			dest.Write(crlf) // data end
			break
		}else{
			// chunk
			dest.Write(chunkSize) // size
			dest.Write(crlf) // size end
			dest.Write(data) // data
			dest.Write(crlf) // data end
		}
	}

	dest.Write(lastChunkAndChunkedBodyEnd) // last-chunk + Chunked-Body end
}