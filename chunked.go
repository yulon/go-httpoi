package httpoi

import (
	"io"
	"strconv"
)

var lastChunkAndChunkedBodyEnd = []byte("0\r\n\r\n")
const chunkDefaultSize = 5120
var chunkDefaultSizeHex = []byte(strconv.FormatUint(uint64(chunkDefaultSize), 16))

func ChunkedEncode(dest io.Writer, src io.Reader) {
	for {
		data := make([]byte, chunkDefaultSize)
		leng, err := src.Read(data)

		if err != nil {
			break
		}

		if leng < chunkDefaultSize {
			// (last-1)-chunk
			dest.Write([]byte(strconv.FormatUint(uint64(leng), 16))) // size
			dest.Write(crlf) // size end
			dest.Write(data[:leng]) // data
			dest.Write(crlf) // data end
			break
		}else{
			// chunk
			dest.Write(chunkDefaultSizeHex) // size
			dest.Write(crlf) // size end
			dest.Write(data) // data
			dest.Write(crlf) // data end
		}
	}

	dest.Write(lastChunkAndChunkedBodyEnd) // last-chunk + Chunked-Body end
}