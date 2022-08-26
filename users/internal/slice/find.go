package slice

func Contains[T comparable](slice []T, el T) bool {
	for i := range slice {
		if slice[i] == el {
			return true
		}
	}

	return false
}
