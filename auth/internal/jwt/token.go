package jwt

import (
	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/lestrrat-go/jwx/jwt"
)

func (m *ManagerImpl) ParseRawToken(accessToken string) (*claims.Subject, error) {
	token, err := jwt.Parse([]byte(accessToken), jwt.WithKeySet(m.pubSet))
	if err != nil {
		return nil, err
	}

	return claims.GetTokenSubject(token)
}
