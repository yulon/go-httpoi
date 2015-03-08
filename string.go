package httpoi

func trim(s string) string {
	leng := len(s)
	for i := 0; i < leng; i++ {
		if s[i] != ' ' && s[i] != '\t' {
			for y := leng - 1; y >= 0; y-- {
				if s[y] != ' ' && s[y] != '\t' {
					return s[i:y+1]
				}
			}
		}
	}
	return ""
}