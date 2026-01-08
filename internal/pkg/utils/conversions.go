package utils

import (
	"net/url"
	"reflect"
	"strconv"
	"time"
)

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

// TimePtr retorna um ponteiro para time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// TimePtrFromString cria um ponteiro para time.Time a partir de string
// Suporta múltiplos formatos common
func TimePtrFromString(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}

	// Tenta diferentes formatos
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t
		}
	}

	return nil
}

func DefaultString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func DefaultInt(value int, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

func BoolPtr(b bool) *bool {
	return &b
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func ParseFloatRange(query url.Values, minKey, maxKey string, minDest, maxDest any) {
	setFloatLike(query.Get(minKey), minDest, parseFloat)
	setFloatLike(query.Get(maxKey), maxDest, parseFloat)
}

func parseFloat(s string) (any, bool) {
	if s == "" {
		return nil, false
	}
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v, true
	}
	// se parse falhar, não atribui nada
	return nil, false
}

func setFloatLike(value string, dest any, parser func(string) (any, bool)) {
	if dest == nil {
		return
	}
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Pointer {
		return
	}
	inner := rv.Elem() // espera-se um pointer (ex: **float64)
	if !inner.IsValid() || inner.Kind() != reflect.Pointer {
		return
	}
	// se valor vazio, não altera (mantém nil)
	if value == "" {
		return
	}
	parsed, ok := parser(value)
	if !ok {
		return
	}

	elemType := inner.Type().Elem()
	switch elemType.Kind() {
	case reflect.Float64:
		newVal := reflect.New(elemType) // *float64
		newVal.Elem().SetFloat(parsed.(float64))
		inner.Set(newVal)
	case reflect.String:
		newVal := reflect.New(elemType) // *string
		newVal.Elem().SetString(value)
		inner.Set(newVal)
	default:
		// tipo não suportado -> nada
	}
}

// ---------------- Int Range ----------------

// ParseIntRange lê dois parâmetros de query (minKey, maxKey) e atribui aos destinos fornecidos.
// Destinos aceitos (passar por referência): **int, **string
func ParseIntRange(query url.Values, minKey, maxKey string, minDest, maxDest any) {
	setIntLike(query.Get(minKey), minDest)
	setIntLike(query.Get(maxKey), maxDest)
}

func setIntLike(value string, dest any) {
	if dest == nil {
		return
	}
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Pointer {
		return
	}
	inner := rv.Elem()
	if !inner.IsValid() || inner.Kind() != reflect.Pointer {
		return
	}
	if value == "" {
		return
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return
	}
	elemType := inner.Type().Elem()
	switch elemType.Kind() {
	case reflect.Int:
		newVal := reflect.New(elemType) // *int
		newVal.Elem().SetInt(int64(parsed))
		inner.Set(newVal)
	case reflect.String:
		newVal := reflect.New(elemType) // *string
		newVal.Elem().SetString(value)
		inner.Set(newVal)
	default:
		// não suportado
	}
}

// ---------------- Time Range ----------------

// ParseTimeRange lê dois parâmetros de query (minKey, maxKey) e atribui aos destinos fornecidos.
// Destinos aceitos (passar por referência): **time.Time, **string
// Aceita vários formatos de data: "2006-01-02", "2006-01-02 15:04:05", RFC3339
func ParseTimeRange(query url.Values, minKey, maxKey string, minDest, maxDest any) {
	setTimeLike(query.Get(minKey), minDest)
	setTimeLike(query.Get(maxKey), maxDest)
}

var timeLayouts = []string{
	"2006-01-02 15:04:05",
	"2006-01-02",
	time.RFC3339,
	"2006-01-02T15:04:05", // variação sem timezone
}

func tryParseTime(s string) (time.Time, bool) {
	for _, l := range timeLayouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func setTimeLike(value string, dest any) {
	if dest == nil {
		return
	}
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Pointer {
		return
	}
	inner := rv.Elem()
	if !inner.IsValid() || inner.Kind() != reflect.Pointer {
		return
	}
	if value == "" {
		return
	}
	t, ok := tryParseTime(value)
	if !ok {
		return
	}
	elemType := inner.Type().Elem()
	switch elemType {
	case reflect.TypeOf(time.Time{}):
		newVal := reflect.New(elemType) // *time.Time
		newVal.Elem().Set(reflect.ValueOf(t))
		inner.Set(newVal)
	case reflect.TypeOf(""):
		newStr := reflect.New(elemType) // *string
		newStr.Elem().SetString(value)
		inner.Set(newStr)
	default:
		// não suportado
	}
}
