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
	Headers Headers
}

type RequestR struct{
	*RequestHeader
	io.Reader
}

func ReadRequest(conn net.Conn) (*RequestR, error) {
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
					if rqLineNowToken != 2 {
						return nil, errors.New("Request Line token count mismatch")
					}

					rqLine.HTTPVersion = buf.String()
					buf.Reset()

					if rqLine.HTTPVersion != "HTTP/1.1" {
						return nil, errors.New("Protocol is not supported")
					}

					rqLine.Method = rqLineTokens[0]
					rqLine.URI = rqLineTokens[1]

					// Read Header Fields
					rqHeaders := Headers{}
					rqHeadersName := ""
					rqHeadersColon := false
					for {
						conn.Read(b)
						switch b[0] {
							case ':':
								if !rqHeadersColon {
									rqHeadersColon = true
									rqHeadersName = buf.String()
									buf.Reset()
								}else{
									buf.Write(b)
								}

							case '\r':
								conn.Read(b)
								if b[0] == '\n' {
									if rqHeadersColon {
										rqHeaders[rqHeadersName] = trim(buf.String())
										buf.Reset()
										rqHeadersName = ""
										rqHeadersColon = false
									}else{
										rq := &RequestR{
											RequestHeader: &RequestHeader{
												RequestLine: rqLine,
												Headers: rqHeaders,
											},
											Reader: conn,
										}
										return rq, nil
									}
								}else{
									buf.Write(b)
								}

							default:
								buf.Write(b)
						}
					}
					return nil, errors.New("Not Header Fields")
				}else{
					buf.Write(b)
				}

			default:
				buf.Write(b)
		}
	}
	return nil, errors.New("Not Request")
}