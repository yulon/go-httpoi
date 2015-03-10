package httpoi

func ParseHeaderField(line []byte) (string, string) {
	for i := 0; i < len(line); i++ {
		if line[i] == ':' {
			return string(line[:i]), trim(string(line[i+1:]))
		}
	}
	return "", ""
}