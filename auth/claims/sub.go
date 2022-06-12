package claims

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwt"
)

type Subject struct {
	UID     string `json:"uid"`
	Expires int64  `json:"exp"`
}

func (s *Subject) UpdateExpiry(exp int64) {
	s.Expires = time.Now().Unix() + exp
}

func GetTokenSubject(token jwt.Token) (*Subject, error) {
	const (
		uidParamName = "uid"
	)

	uid, ok := token.Get(uidParamName)
	if !ok {
		return nil, fmt.Errorf("no %s param", uidParamName)
	}

	uidString, ok := uid.(string)
	if !ok {
		return nil, fmt.Errorf("invalid type of %s: %t", uidParamName, uid)
	}

	return &Subject{
		UID:     uidString,
		Expires: token.Expiration().Unix(),
	}, nil
}
