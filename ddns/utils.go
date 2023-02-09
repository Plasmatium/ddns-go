package main

func Ref[T any](v T) *T {
	return &v
}
