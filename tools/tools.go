package tools

import (
	"os"
	"strings"
)

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func StrFirstToUpper(str string) string {
	var res string
	var up bool
	for k, ch := range str {
		if ch == '_' {
			up = true
			continue
		}
		if k == 0 || up {
			res += strings.ToUpper(string(ch))
			up = false
			continue
		}

		res += string(ch)
	}
	return res
}
