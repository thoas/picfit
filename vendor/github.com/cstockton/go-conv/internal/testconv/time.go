package testconv

import (
	"testing"
	"time"
)

func RunTimeTests(t *testing.T, fn func(interface{}) (time.Time, error)) {
	RunTest(t, TimeKind, func(v interface{}) (interface{}, error) {
		return fn(v)
	})
}

type testTimeConverter time.Time

func (t testTimeConverter) Time() (time.Time, error) {
	return time.Time(t).Add(time.Minute), nil
}

func init() {
	var (
		emptyTime time.Time
		t2006     = time.Date(2006, time.January, 2, 15, 4, 5, 999999999, time.UTC)
	)

	// basic
	assert(time.Time{}, time.Time{})
	assert(new(time.Time), time.Time{})
	assert(t2006, t2006)

	// strings
	fmts := []string{
		"02 Jan 06 15:04:05",
		"2 Jan 2006 15:04:05",
		"2 Jan 2006 15:04:05 -0700 (UTC)",
		"02 Jan 2006 15:04 UTC",
		"02 Jan 2006 15:04:05 UTC",
		"02 Jan 2006 15:04:05 -0700 (UTC)",
		"Mon, 2 Jan  15:04:05 UTC 2006",
		"Mon, 2 Jan 15:04:05 UTC 2006",
		"Mon, 02 Jan 2006 15:04:05",
		"Mon, 02 Jan 2006 15:04:05 (UTC)",
		"Mon, 2 Jan 2006 15:04:05",
	}
	for _, s := range fmts {
		assert(s, TimeExp{Moment: t2006.Truncate(time.Minute), Truncate: time.Minute})
		assert(testStringConverter(s), TimeExp{Moment: t2006.Truncate(time.Minute), Truncate: time.Minute})
	}

	// underlying
	type ulyTime time.Time
	assert(ulyTime(t2006), t2006)
	assert(ulyTime(t2006), t2006)

	// implements converter
	assert(testTimeConverter(t2006), t2006.Add(time.Minute))

	// embedded time
	type embedTime struct{ time.Time }
	assert(embedTime{t2006}, t2006)

	// errors
	assert(nil, experr(emptyTime, `cannot convert <nil> (type <nil>) to time.Time`))
	assert("foo", experr(emptyTime, `cannot convert "foo" (type string) to time.Time`))
	assert("tooLong", experr(
		emptyTime, `cannot convert "tooLong" (type string) to time.Time`))
	assert(struct{}{}, experr(
		emptyTime, `cannot convert struct {}{} (type struct {}) to `))
	assert([]string{"1s"}, experr(
		emptyTime, `cannot convert []string{"1s"} (type []string) to `))
	assert([]string{}, experr(
		emptyTime, `cannot convert []string{} (type []string) to `))
}
