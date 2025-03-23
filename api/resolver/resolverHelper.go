package resolver

func ptr[T interface{}](in T) *T {
	return &in
}
