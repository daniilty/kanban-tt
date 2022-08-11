package core

type TTL struct {
	Default int
	List    []int
}

func (s *ServiceImpl) GetDefaultTTL() int64 {
	return defaultTTL
}

func (s *ServiceImpl) GetTTLs() []int64 {
	listCopy := make([]int64, 0, len(ttlList))
	copy(listCopy, ttlList)

	return listCopy
}

const (
	defaultTTL = 365
)

var (
	ttlList = []int64{1, 2, 3, 5, 10, 30, 90, 180, 365}
)
