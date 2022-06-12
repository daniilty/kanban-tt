package jwt

func (m *ManagerImpl) Refresh(accessToken string) (string, error) {
	sub, err := m.ParseRawToken(accessToken)
	if err != nil {
		return "", err
	}

	return m.Generate(sub)
}
