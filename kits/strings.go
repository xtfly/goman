package kits

import (
	"regexp"
	"strings"

	"github.com/Unknwon/com"
)

func IsDigit(str string) bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(str)
}

func IsChineseLetterNumUnline(str string) bool {
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}_a-zA-Z0-9]+$`).MatchString(str)
}

func IsLetterNumUnline(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(str)
}

func IsChinese(str string) bool {
	return regexp.MustCompile(`^[\x{4e00}-\x{9fa5}]+$`).MatchString(str)
}

func IsEmail(str string) bool {
	//return regexp.MustCompile(`^([a-z0-9\+_\-]+)(\.[a-z0-9\+_\-]+)*@([a-z0-9\-]+\.)+[a-z]{2,6}$`).MatchString(str)
	return com.IsEmail(str)
}

// 是否包含敏感词
func SensitiveWordExists(str string, sw string) bool {
	if len(str) == 0 || len(sw) == 0 {
		return false
	}

	ws := strings.Split(sw, "\n")
	for _, w := range ws {
		w = strings.TrimSpace(w)
		if len(w) == 0 {
			continue
		}

		if strings.HasPrefix(w, "{") && strings.HasSuffix(w, "}") {
			if regexp.MustCompile(`^` + w[1:len(w)-1] + `$`).MatchString(str) {
				return true
			}
		}

		if strings.Contains(str, w) {
			return true
		}
	}

	return false
}

func SensitiveWordExistsV2(strs []string, sw string) bool {
	for _, s := range strs {
		if SensitiveWordExists(s, sw) {
			return true
		}
	}

	return false
}

func IfEmpty(a, b string) string {
	if a == "" {
		return b
	}
	return a
}
