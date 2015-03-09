package httpoi

var crlf = []byte("\r\n")

func concat(d ...[]byte) []byte {
	num := len(d)
	leng := 0
	for i := 0; i < num; i++ {
		leng += len(d[i])
	}
	buf := make([]byte, leng, leng)
	leng = 0
	for i := 0; i < num; i++ {
		copy(buf[leng:], d[i])
		leng += len(d[i])
	}
	return buf
}