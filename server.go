package httpoi

import (
	"net"
	//"fmt"
	"time"
	"errors"
	"runtime"
)

type SeverHandler func(*ResponseW, *RequestR)

func Sever(laddr string, h SeverHandler) error {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return errors.New("Failed to listen")
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return errors.New("Failed to accept")
		}
		go saw(conn, h)
	}
	return nil
}

var lang = runtime.Version()

func saw(conn net.Conn, h SeverHandler) {
	rq, err := ReadRequest(conn)
	if err == nil {
		rs := &ResponseW{
			line: &StatusLine{
				HTTPVersion: rq.HTTPVersion,
			},
			Headers: map[string]string{
				"Date": time.Now().Format(time.RFC1123),
				"Server": "HTTPOI",
				"X-Powered-By": lang,
			},
			wc: &writeNopCloser{conn},
		}
		h(rs, rq)
		rs.wc.Close()
	}
	conn.Close()
}