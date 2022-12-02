package utils

func Has[T comparable](o T, list []T) bool {
	for _, r := range list {
		if r == o {
			return true
		}
	}
	return false
}

func Map[T any, R any](list []T, f func(t T) R) []R {
	r := make([]R, 0)
	for _, t := range list {
		r = append(r, f(t))
	}
	return r
}
