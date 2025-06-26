package utils

func Int64OrNil(ptr *int64) interface{} {
	if ptr == nil {
		return nil
	}
	return *ptr
}
