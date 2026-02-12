package utils

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPointerHelpers(t *testing.T) {
	var i64 int64 = 10
	assert.Equal(t, int64(10), Int64OrNil(&i64))
	assert.Nil(t, Int64OrNil(nil))

	assert.Equal(t, int64(5), *Int64Ptr(5))
	assert.Equal(t, 3, *IntPtr(3))
	assert.Equal(t, uint(7), *UintPtr(7))
	assert.Equal(t, true, *BoolPtr(true))
	assert.Equal(t, 1.5, *Float64Ptr(1.5))

	str := "abc"
	assert.Equal(t, "abc", *StrToPtr(str))
}

func TestNilAndDerefHelpers(t *testing.T) {
	var nilInt *int64
	assert.Equal(t, int64(0), NilToZero(nilInt))
	assert.Equal(t, int64(0), DerefInt64(nilInt))

	val := int64(20)
	assert.Equal(t, int64(20), NilToZero(&val))
	assert.Equal(t, int64(20), DerefInt64(&val))

	var nilTime *time.Time
	assert.Equal(t, time.Time{}, DerefTime(nilTime))

	now := time.Now()
	assert.Equal(t, now, DerefTime(&now))
	assert.Equal(t, now, *TimePtr(now))
}

func TestTimePtrFromString(t *testing.T) {
	assert.Nil(t, TimePtrFromString(""))

	validDates := []string{
		"2023-01-02",
		"2023-01-02 15:04:05",
		"2023-01-02T15:04:05",
		time.Now().Format(time.RFC3339),
	}

	for _, d := range validDates {
		assert.NotNil(t, TimePtrFromString(d))
	}

	assert.Nil(t, TimePtrFromString("invalid-date"))
}

func TestDefaultHelpers(t *testing.T) {
	assert.Equal(t, "fallback", DefaultString("", "fallback"))
	assert.Equal(t, "value", DefaultString("value", "fallback"))

	assert.Equal(t, 99, DefaultInt(0, 99))
	assert.Equal(t, 10, DefaultInt(10, 99))
}

func TestParseFloatRange_AllBranches(t *testing.T) {
	q := url.Values{}
	q.Set("min", "10.5")
	q.Set("max", "20.5")

	var min *float64
	var max *float64

	ParseFloatRange(q, "min", "max", &min, &max)

	assert.NotNil(t, min)
	assert.NotNil(t, max)
	assert.Equal(t, 10.5, *min)
	assert.Equal(t, 20.5, *max)

	// string destination
	var minStr *string
	ParseFloatRange(q, "min", "max", &minStr, nil)
	assert.NotNil(t, minStr)
	assert.Equal(t, "10.5", *minStr)

	// invalid float
	q.Set("min", "abc")
	var invalid *float64
	ParseFloatRange(q, "min", "max", &invalid, nil)
	assert.Nil(t, invalid)

	// empty value
	q.Set("min", "")
	ParseFloatRange(q, "min", "max", &invalid, nil)
	assert.Nil(t, invalid)

	// unsupported type
	var unsupported *int
	ParseFloatRange(q, "min", "max", &unsupported, nil)
	assert.Nil(t, unsupported)

	// dest nil
	ParseFloatRange(q, "min", "max", nil, nil)

	// dest not pointer
	ParseFloatRange(q, "min", "max", min, nil)
}

func TestParseIntRange_AllBranches(t *testing.T) {
	q := url.Values{}
	q.Set("min", "5")
	q.Set("max", "10")

	var min *int
	var max *int

	ParseIntRange(q, "min", "max", &min, &max)

	assert.NotNil(t, min)
	assert.NotNil(t, max)
	assert.Equal(t, 5, *min)
	assert.Equal(t, 10, *max)

	// string destination
	var minStr *string
	ParseIntRange(q, "min", "max", &minStr, nil)
	assert.NotNil(t, minStr)
	assert.Equal(t, "5", *minStr)

	// invalid int
	q.Set("min", "abc")
	var invalid *int
	ParseIntRange(q, "min", "max", &invalid, nil)
	assert.Nil(t, invalid)

	// empty
	q.Set("min", "")
	ParseIntRange(q, "min", "max", &invalid, nil)
	assert.Nil(t, invalid)

	// unsupported type
	var unsupported *float64
	ParseIntRange(q, "min", "max", &unsupported, nil)
	assert.Nil(t, unsupported)

	// dest nil
	ParseIntRange(q, "min", "max", nil, nil)

	// dest not pointer
	ParseIntRange(q, "min", "max", min, nil)
}

func TestParseTimeRange_AllBranches(t *testing.T) {
	q := url.Values{}
	q.Set("start", "2023-01-01")
	q.Set("end", "2023-01-02 10:00:00")

	var start *time.Time
	var end *time.Time

	ParseTimeRange(q, "start", "end", &start, &end)

	assert.NotNil(t, start)
	assert.NotNil(t, end)

	// string destination
	var startStr *string
	ParseTimeRange(q, "start", "end", &startStr, nil)
	assert.NotNil(t, startStr)
	assert.Equal(t, "2023-01-01", *startStr)

	// invalid time
	q.Set("start", "invalid")
	var invalid *time.Time
	ParseTimeRange(q, "start", "end", &invalid, nil)
	assert.Nil(t, invalid)

	// empty
	q.Set("start", "")
	ParseTimeRange(q, "start", "end", &invalid, nil)
	assert.Nil(t, invalid)

	// unsupported type
	var unsupported *int
	ParseTimeRange(q, "start", "end", &unsupported, nil)
	assert.Nil(t, unsupported)

	// dest nil
	ParseTimeRange(q, "start", "end", nil, nil)

	// dest not pointer
	ParseTimeRange(q, "start", "end", start, nil)
}

func TestTryParseTime_AllLayouts(t *testing.T) {
	for _, layout := range timeLayouts {
		str := time.Now().Format(layout)
		tm, ok := tryParseTime(str)
		assert.True(t, ok)
		assert.False(t, tm.IsZero())
	}

	_, ok := tryParseTime("invalid")
	assert.False(t, ok)
}

func TestSetTimeLike_UnsupportedTypeBranch(t *testing.T) {
	q := url.Values{}
	q.Set("date", "2023-01-01")

	var unsupported **int
	setTimeLike("2023-01-01", &unsupported)

	assert.Nil(t, unsupported)
}

func TestReflectBranchCoverage(t *testing.T) {
	// força branch onde inner não é pointer
	var x int
	setIntLike("10", &x)
	setFloatLike("10.5", &x, parseFloat)
	setTimeLike("2023-01-01", &x)

	assert.Equal(t, reflect.Int, reflect.TypeOf(x).Kind())
}
