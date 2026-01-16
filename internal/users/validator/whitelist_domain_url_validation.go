package validator

import "strings"

func IsSafeRedirectURL(returnURL string, allowedPrefix string) bool {
	if returnURL == "" {
		return false
	}

	if strings.HasPrefix(returnURL, "http://") ||
		strings.HasPrefix(returnURL, "https://") {
		return false
	}

	if strings.HasPrefix(returnURL, "//") {
		return false
	}

	if !strings.HasPrefix(returnURL, allowedPrefix) {
		return false
	}

	return true
}
