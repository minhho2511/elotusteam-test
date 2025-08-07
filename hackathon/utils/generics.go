package utils

func WithPointer[T any](p T) *T {
	return &p
}

func Contains[T comparable](array []T, el T) bool {
	for _, a := range array {
		if a == el {
			return true
		}
	}
	return false
}

func Reverse[T any](s []T) []T {
	n := len(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
