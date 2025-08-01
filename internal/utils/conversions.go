package utils

func Int64OrNil(ptr *int64) interface{} {
	if ptr == nil {
		return nil
	}
	return *ptr
}

func ToPointer[T any](v T) *T {
	return &v
}
