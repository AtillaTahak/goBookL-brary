package url

import (
	"net/url"
	"strings"
)

func Canonicalize(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.Path = strings.TrimSuffix(parsed.Path, "/")
	return parsed.String(), nil
}

func Redirect(raw string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	host := strings.ToLower(parsed.Host)
	if !strings.HasPrefix(host, "www.") {
		host = "www." + host
	}
	parsed.Host = host

	parsed.Scheme = "https"

	parsed.Path = strings.ToLower(parsed.Path)

	return parsed.String(), nil
}

func CleanURL(raw, operation string) (string, error) {
	var err error
	result := raw
	switch operation {
	case "canonical":
		result, err = Canonicalize(raw)
	case "redirection":
		result, err = Redirect(raw)
	case "all":
		result, err = Canonicalize(raw)
		if err != nil {
			return "", err
		}
		result, err = Redirect(result)
	}
	return result, err
}
