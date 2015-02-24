package httpoi

import (
	"time"
	"strconv"
)

func TimeToUnixString(srctime time.Time) string {
	return strconv.Itoa(int(srctime.Unix()))
}

func TimeToHttpString(srctime time.Time) string {
	return srctime.Format(time.RFC1123)
}

func hash(data []byte) int {
	var h int
	for i := 0; i < len(data); i++ {
		h = h * 31 + int(data[i])
	}
	return h
}