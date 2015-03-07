package httpoi

import (
	"io"
)

type writeNopCloser struct{
	io.Writer
}

func (writeNopCloser) Close() error { return nil }