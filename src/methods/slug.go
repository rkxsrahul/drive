package methods

import (
	"strings"
)

func SlugOfString(name string) string {
	// replace special characters with space
	str := []byte(name)
	for i := 0; i < len(str); i++ {
		if (str[i] > 47 && str[i] < 58) || (str[i] > 64 && str[i] < 91) || (str[i] > 96 && str[i] < 123) {
			continue
		} else {
			str[i] = 32
		}
	}
	// remove consecutive spaces
	name = string(str)
	name = strings.Join(strings.Fields(name), " ")
	// create slug of string by replacing space with -
	slug := ""
	for i := 0; i < len(name); i++ {
		chr := name[i]
		if chr == ' ' {
			chr = '-'
		}
		slug += string(chr)
	}
	return strings.ToLower(slug)
}
