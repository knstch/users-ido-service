package validator

import (
	"net/url"
	"path"
	"strings"
)

func IsSafeRedirectURL(returnURL string, allowedPrefix string) bool {
	if returnURL == "" {
		return false
	}

	if strings.HasPrefix(returnURL, "//") {
		return false
	}

	// Allow relative paths (starting with /)
	if strings.HasPrefix(returnURL, "/") {
		return true
	}

	// Parse absolute URL
	u, err := url.Parse(returnURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Parse allowed prefix to validate domain
	base, err := url.Parse(allowedPrefix)
	if err != nil || base.Scheme == "" || base.Host == "" {
		return false
	}
	// Must match scheme and host
	if u.Scheme != base.Scheme || u.Host != base.Host {
		return false
	}

	// If base has a path, the target path must be under it
	if base.Path != "" && base.Path != "/" {
		basePath := path.Clean(base.Path)
		targetPath := path.Clean(u.Path)
		if !strings.HasPrefix(targetPath+"/", basePath+"/") {
			return false
		}
	}

	return true
}
