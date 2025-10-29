package utils

import "time"

func Int64OrNil(ptr *int64) any {
	if ptr == nil {
		return nil
	}
	return *ptr
}

func Int64Ptr(i int64) *int64 {
	return &i
}

func IntPtr(i int) *int {
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

func UintPtr(v uint) *uint {
	return &v
}

func DerefInt64(p *int64) int64 {
	if p != nil {
		return *p
	}
	return 0
}

func DerefTime(p *time.Time) time.Time {
	if p != nil {
		return *p
	}
	return time.Time{}
}

func DefaultString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
