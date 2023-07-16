package swrcache

func contains[T comparable](s []T, target T) bool {
	for _, elt := range s {
		if elt == target {
			return true
		}
	}
	return false
}
