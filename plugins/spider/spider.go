package spider

import (
	"net/http"
	"strings"

	"gopkg.in/macaron.v1"
)

var (
	searchengineBot []string = []string{
		"googlebot",
		"mediapartners-google",
		"msnbot",
		"yodaobot",
		"sosospider+",
		"yahoo! slurp",
		"spider",
	}
)

// 检查是否为搜索引擎爬虫
func isSpider(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)
	for _, v := range searchengineBot {
		if strings.Contains(userAgent, v) {
			return true
		}
	}
	return false
}

func SpiderFunc() macaron.Handler {
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		userAgent := req.Header.Get("User-Agent")
		if isSpider(userAgent) {
			http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}
}
