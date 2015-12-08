package account

import (
	"regexp"
	"strings"
)

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
