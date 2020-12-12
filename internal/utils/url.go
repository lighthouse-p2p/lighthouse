package utils

import "strings"

// TranslateURL translates all non local URLs to use HTTPS/WSS
func TranslateURL(url string) string {
	if strings.Contains(url, "192.168") || strings.Contains(url, "10.") || strings.Contains(url, "localhost") {
		return url
	}

	if !strings.HasPrefix(url, "https") || !strings.HasPrefix(url, "wss") {
		if strings.HasPrefix(url, "http") {
			return strings.Replace(url, "http", "https", 1)
		}

		if strings.HasPrefix(url, "ws") {
			return strings.Replace(url, "ws", "wss", 1)
		}
	}

	return url
}
