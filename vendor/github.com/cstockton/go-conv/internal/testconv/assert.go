package testconv

import (
	"errors"
	"fmt"
	"math"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/cstockton/go-conv/internal/refutil"
)

var (
	assertions      Assertions
	assertionsIndex = make(AssertionsIndex)
)

const (
	TypeKind reflect.Kind = reflect.UnsafePointer + iota
	DurationKind
	TimeKind
	SliceMask reflect.Kind = (1 << 8)
	MapMask                = (1 << 16)
	ChanMask               = (1 << 24)
)

var refKindNames = map[reflect.Kind]string{
	DurationKind: "time.Duration",
	TimeKind:     "time.Time",
}

func reflectHasMask(k, mask reflect.Kind) bool {
	return (k & mask) == mask
}

func reflectKindStr(k reflect.Kind) string {
	if reflectHasMask(k, ChanMask|SliceMask|MapMask) {
		return refKindNames[k]
	}
	return k.String()
}

// returns the reflect.Kind with the additional masks / types for conversions.
func convKindStr(v interface{}) string { return reflectKindStr(convKind(v)) }
func convKind(v interface{}) reflect.Kind {
	switch T := refutil.Indirect(v).(type) {
	case reflect.Kind:
		return T
	case Expecter:
		return T.Kind()
	case time.Time:
		return TimeKind
	case time.Duration:
		return DurationKind
	default:
		kind := reflect.TypeOf(T).Kind()
		if reflect.Slice == kind || kind == reflect.Array {
			kind |= SliceMask
		} else if kind == reflect.Map {
			kind |= MapMask
		} else if kind == reflect.Chan {
			kind |= SliceMask
		}
		return kind
	}
}

// Expecter defines a type of conversion for v and will return an error if the
// value is unexpected.
type Expecter interface {

	// Exp is given the result of conversion and an error if one occurred. The
	// type of conversion operation should be the type of Kind().
	Expect(got interface{}, err error) error

	// Value is the expected result.
	Value() interface{}

	// Kind returns the type of conversion that should be given to Exepct().
	Kind() reflect.Kind
}

type AssertionsLookup struct {
	assertion *Assertion
	expecter  Expecter
}

// AssertionsIndex is map of assertions by the Expecter kind.
type AssertionsIndex map[reflect.Kind][]AssertionsLookup

// Assertions is slice of assertions.
type Assertions []*Assertion

// EachOf will visit each assertion that contains an Expecter for reflect.Kind.
func (a Assertions) EachOf(k reflect.Kind, f func(a *Assertion, e Expecter)) int {
	found := assertionsIndex[k]
	for _, v := range found {
		f(v.assertion, v.expecter)
	}
	return len(found)
}

// Assertion represents a set of expected values from a Converter. It's the
// general implementation of asserter. It will compare each Expects to the
// associated Converter.T(). It returns an error for each failed conversion
// or an empty error slice if no failures occurred.
type Assertion struct {

	// From is the underlying value to be used in conversions.
	From interface{}

	// Name of this assertion
	Name string

	// Code to construct Interface
	Code string

	// Type of Interface field using fmt %#T.
	Type string

	// File is the file the assertion is defined in.
	File string

	// Line is the line number the assertion is defined in.
	Line int

	// Description for this assertion.
	Desc string

	// Expects contains a list of values this Assertion expects to see as a result
	// of calling the conversion function for the associated type. For example
	// appending a single int64 to this slice would mean it is expected that:
	//   conv.Int64(assertion.From) == Expect[0].(int64)
	// When appending multiple values of the same type, the last one appended is
	// used.
	Exps []Expecter
}

func (a Assertion) String() string {
	return fmt.Sprintf("[assertion %v:%d] from `%v` (%[3]T)",
		a.File, a.Line, a.From)
}

// assert will create an assertion type for the first argument given. Args is
// consumed into the Expects slice, with slices of interfaces being flattened.
func assert(value interface{}, args ...interface{}) {
	a := &Assertion{
		From: value,
		Type: typename(value),
	}
	split := strings.SplitN(a.Type, ".", 2)
	a.Name = strings.Title(split[len(split)-1])
	a.Code = fmt.Sprintf("%#v", value)
	_, a.File, a.Line, _ = runtime.Caller(1)
	a.File = path.Base(a.File)

	var slurp func(interface{})
	slurp = func(value interface{}) {
		switch T := value.(type) {
		case []interface{}:
			for _, t := range T {
				slurp(t)
			}
		case []Expecter:
			for _, t := range T {
				slurp(t)
			}
		case Expecter:
			a.Exps = append(a.Exps, T)
		default:
			a.Exps = append(a.Exps, Exp{value})
		}
	}
	for _, arg := range args {
		slurp(arg)
	}
	for _, exp := range a.Exps {
		k := exp.Kind()
		assertionsIndex[k] = append(assertionsIndex[k], AssertionsLookup{a, exp})
	}
	assertions = append(assertions, a)
}

// Exp is used for ducmentation purposes.
type Exp struct {
	Want interface{}
}

func (e Exp) Kind() reflect.Kind {
	return convKind(e.Want)
}

func (e Exp) Value() interface{} {
	return e.Want
}

func (e Exp) Expect(got interface{}, err error) error {
	if err == nil && !reflect.DeepEqual(got, e.Want) {
		return fmt.Errorf("exp: (%T) %[1]v\ngot: (%T) %[2]v", e.Want, got)
	}
	return err
}

func (e Exp) String() string {
	return fmt.Sprintf("to %v (type %[1]T)", e.Want)
}

// ErrorExp ensures a given conversion failed.
type ErrorExp struct {
	Exp
	ErrStr string
}

func experr(want interface{}, contains string) Expecter {
	return ErrorExp{ErrStr: contains, Exp: Exp{Want: want}}
}

func (e ErrorExp) Expect(got interface{}, err error) error {
	if err != nil {
		if len(e.ErrStr) == 0 {
			return err
		}
		if !strings.Contains(err.Error(), e.ErrStr) {
			return fmt.Errorf("error did not match:\n  exp: %v\n  got: %v", e.ErrStr, err)
		}
	} else if len(e.ErrStr) > 0 {
		return errors.New("expected non-nil err")
	}
	return nil
}

// Float64Exp asserts that (converter.Float64() - Exp) < epsilon64
type Float64Exp struct {
	Want float64
}

const (
	epsilon64 = float64(.00000000000000001)
)

func (e Float64Exp) Kind() reflect.Kind { return reflect.Float64 }
func (e Float64Exp) Expect(got interface{}, err error) error {
	val := got.(float64)
	abs := math.Abs(val - e.Want)
	if abs < epsilon64 {
		return err
	}
	return fmt.Errorf("%#v.assert(%#v): abs value %v exceeded epsilon %v",
		e, val, abs, epsilon64)
}

func (e Float64Exp) Value() interface{} {
	return e.Want
}

// Float32Exp asserts that (converter.Float32() - Exp) < epsilon64
type Float32Exp struct {
	Want float32
}

func (e Float32Exp) Kind() reflect.Kind { return reflect.Float32 }
func (e Float32Exp) Expect(got interface{}, err error) error {
	val := got.(float32)
	abs := math.Abs(float64(val - e.Want))
	if abs < epsilon64 {
		return err
	}
	return fmt.Errorf("%#v.assert(%#v): abs value %v exceeded epsilon %v",
		e, val, abs, epsilon64)
}

func (e Float32Exp) Value() interface{} {
	return e.Want
}

// TimeExp helps validate time.Time() conversions, specifically because under
// some conversions time.Now() may be used. It will check that the difference
// between the Moment is the same as the given value after truncation and
// rounding (if either is set) is identical.
type TimeExp struct {
	Moment   time.Time
	Offset   time.Duration
	Round    time.Duration
	Truncate time.Duration
}

func (e TimeExp) String() string     { return fmt.Sprintf("%v", e.Moment) }
func (e TimeExp) Kind() reflect.Kind { return TimeKind }
func (e TimeExp) Expect(got interface{}, err error) error {
	val := got.(time.Time).Add(e.Offset)
	if e.Round != 0 {
		val = val.Round(e.Round)
	}
	if e.Truncate != 0 {
		val = val.Round(e.Truncate)
	}
	if !e.Moment.Equal(val) {
		return fmt.Errorf(
			"times did not match:\n  exp: %v\n  got: %v", e.Moment, val)
	}
	return nil
}

func (e TimeExp) Value() interface{} {
	return e.Moment
}

// DurationExp supports fuzzy duration conversions.
type DurationExp struct {
	Want  time.Duration
	Round time.Duration
}

func (e DurationExp) Kind() reflect.Kind { return DurationKind }
func (e DurationExp) Expect(got interface{}, err error) error {
	d := got.(time.Duration)
	if e.Round != 0 {
		neg := d < 0
		if d < 0 {
			d = -d
		}
		if m := d % e.Round; m+m < e.Round {
			d = d - m
		} else {
			d = d + e.Round - m
		}
		if neg {
			d = -d
		}
	}
	if e.Want != d {
		return fmt.Errorf("%#v: %#v != %#v", e, e.Want, d)
	}
	return nil
}

func (e DurationExp) Value() interface{} {
	return e.Want
}

func typename(value interface{}) (name string) {
	parts := strings.SplitN(fmt.Sprintf("%T", value), ".", 2)

	if len(parts) == 2 {
		name = parts[0] + "." + parts[1]
	} else {
		name = parts[0]
	}
	return
}
