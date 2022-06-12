package claims

import (
	"net/http"

	"github.com/lestrrat-go/jwx/jwt"
)

func ParseHTTPHeader(h http.Header) (*Subject, error) {
	const authorization = "Authorization"

	token, err := jwt.ParseHeader(h, authorization)
	if err != nil {
		return nil, err
	}

	return GetTokenSubject(token)
}
