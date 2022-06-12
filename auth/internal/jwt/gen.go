package jwt

import (
	"encoding/json"

	"github.com/daniilty/kanban-tt/auth/claims"
	"github.com/lestrrat-go/jwx/jws"
)

func (m *ManagerImpl) Generate(sub *claims.Subject) (string, error) {
	sub.UpdateExpiry(m.tokenExp)

	bb, err := json.Marshal(sub)
	if err != nil {
		return "", err
	}

	accessToken, err := jws.Sign(bb, m.alg, m.privateKey)

	return string(accessToken), err
}
