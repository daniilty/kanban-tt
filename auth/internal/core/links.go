package core

import (
	"net/url"
	"path"
)

func generateConfirmLink(base *url.URL, token string) *url.URL {
	urlCopy := *base
	urlCopy.Path = path.Join(urlCopy.Path, token)

	return &urlCopy
}
