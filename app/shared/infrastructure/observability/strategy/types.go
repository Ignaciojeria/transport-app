package strategy

import (
	"net/url"
	"strings"
)

type OpenObserveHttpEndpoint string

func (o OpenObserveHttpEndpoint) GetDNS() string {
	parsedURL, err := url.Parse(string(o))
	if err != nil {
		return ""
	}
	host := parsedURL.Host
	if host == "" {
		parts := strings.SplitN(string(o), "/", 2)
		if len(parts) > 0 {
			host = parts[0]
		}
	}

	return host
}

func (o OpenObserveHttpEndpoint) GetPath() string {
	parsedURL, err := url.Parse(string(o))
	if err != nil {
		return ""
	}
	return parsedURL.Path
}
