package core

func (s *ServiceImpl) RefreshSession(accessToken string) (string, error) {
	return s.jwtManager.Refresh(accessToken)
}

func (s *ServiceImpl) JWKS() []byte {
	return s.jwtManager.JWKS()
}
