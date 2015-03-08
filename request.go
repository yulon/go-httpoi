package httpoi

import (
	"io"
	"errors"
	"bytes"
	"net"
	//"fmt"
)

type RequestLine struct{
	Method string
	URI string
	HTTPVersion string
}

type RequestHeader struct{
	*RequestLine
	HF HeaderFields
}

type RequestR struct{
	*RequestHeader
	io.Reader
}

func ReadRequest(conn net.Conn) (rq *RequestR, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))

	// Read Request Line
	rqLine := &RequestLine{}
	rqLineTokens := make([]string, 2, 2)
	rqLineNowToken := 0
	b := make([]byte, 1, 1)
	for {
		conn.Read(b)
		switch b[0] {
			case ' ':
				rqLineTokens[rqLineNowToken] = buf.String()
				buf.Reset()
				rqLineNowToken++

			case '\r':
				conn.Read(b)
				if b[0] == '\n' {
					if rqLineNowToken == 2 {
						rqLine.HTTPVersion = buf.String()
						buf.Reset()
						if rqLine.HTTPVersion == "HTTP/1.1" {
							rqLine.Method = rqLineTokens[0]
							rqLine.URI = rqLineTokens[1]

							// Read Request Header Fields
							rqHF := HeaderFields{}
							rqHFName := ""
							rqHFColon := false
							for {
								conn.Read(b)
								switch b[0] {
									case ':':
										if !rqHFColon {
											rqHFColon = true
											rqHFName = buf.String()
											buf.Reset()
										}else{
											buf.Write(b)
										}

									case '\r':
										conn.Read(b)
										if b[0] == '\n' {
											if rqHFName != "" {
												rqHF[rqHFName] = trim(buf.String())
												buf.Reset()
												rqHFName = ""
												rqHFColon = false
											}else{
												rq = &RequestR{
													RequestHeader: &RequestHeader{
														RequestLine: rqLine,
														HF: rqHF,
													},
													Reader: conn,
												}
												return
											}
											
										}else{
											buf.Write(b)
										}

									default:
										buf.Write(b)
								}
							}
						}else{
							err = errors.New("Protocol is not supported")
						}
					}else{
						err = errors.New("Request-Line token count mismatch")
					}
					return
				}else{
					buf.Write(b)
				}

			default:
				buf.Write(b)
		}
	}
}