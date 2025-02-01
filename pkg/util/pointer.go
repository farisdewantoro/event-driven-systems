package util

func ToPointer[T interface{}](v T) *T {
	return &v
}
