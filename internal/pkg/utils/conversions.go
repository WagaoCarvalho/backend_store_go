package utils

func Int64OrNil(ptr *int64) any {
	if ptr == nil {
		return nil
	}
	return *ptr
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func StrToPtr[T any](v T) *T {
	return &v
}

func NilToZero(id *int64) int64 {
	if id == nil {
		return 0
	}
	return *id
}
