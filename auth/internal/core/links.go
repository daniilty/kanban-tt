package core

import (
	"net/url"
)

func generateConfirmLink(base *url.URL, token string) *url.URL {
	urlCopy := *base
	urlCopy.RawQuery = "token=" + token

	return &urlCopy
}
