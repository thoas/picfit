package refconv

import (
	"flag"
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/cstockton/go-conv/internal/testconv"
)

// Test converter.
func TestMain(m *testing.M) {
	chkMathIntSize := mathIntSize
	chkMathMaxInt := mathMaxInt
	chkMathMinInt := mathMinInt
	chkMathMaxUint := mathMaxUint
	chkEmptyTime := time.Time{}

	flag.Parse()
	res := m.Run()

	// validate our max (u|)int sizes don't get written to on accident.
	if chkMathIntSize != mathIntSize {
		panic("chkEmptyTime != emptyTime")
	}
	if chkMathMaxInt != mathMaxInt {
		panic("chkMathMaxInt != mathMaxInt")
	}
	if chkMathMinInt != mathMinInt {
		panic("chkMathMinInt != mathMaxInt")
	}
	if chkMathMaxUint != mathMaxUint {
		panic("chkMathMaxUint != mathMaxUint")
	}
	if chkEmptyTime != emptyTime {
		panic("chkEmptyTime != emptyTime")
	}

	os.Exit(res)
}

// Runs common tests to make sure conversion rules are satisfied.
func TestConv(t *testing.T) {
	var c Conv
	testconv.RunBoolTests(t, c.Bool)
	testconv.RunDurationTests(t, c.Duration)
	testconv.RunFloat32Tests(t, c.Float32)
	testconv.RunFloat64Tests(t, c.Float64)
	testconv.RunInferTests(t, c.Infer)
	testconv.RunIntTests(t, c.Int)
	testconv.RunInt8Tests(t, c.Int8)
	testconv.RunInt16Tests(t, c.Int16)
	testconv.RunInt32Tests(t, c.Int32)
	testconv.RunInt64Tests(t, c.Int64)
	testconv.RunStringTests(t, c.String)
	testconv.RunTimeTests(t, c.Time)
	testconv.RunUintTests(t, c.Uint)
	testconv.RunUint8Tests(t, c.Uint8)
	testconv.RunUint16Tests(t, c.Uint16)
	testconv.RunUint32Tests(t, c.Uint32)
	testconv.RunUint64Tests(t, c.Uint64)
}

func TestConvHelpers(t *testing.T) {
	var c Conv
	t.Run("convNumToBool", func(t *testing.T) {
		var val reflect.Value
		if got, ok := c.convNumToBool(0, val); ok || got {
			t.Fatal("expected failure")
		}
	})
	t.Run("convNumToDuration", func(t *testing.T) {
		var val reflect.Value
		if _, ok := c.convNumToDuration(0, val); ok {
			t.Fatal("expected convNumToDuration to return false on invalid kind")
		}
	})
	t.Run("timeFromString", func(t *testing.T) {
		if _, ok := convStringToTime(""); ok {
			t.Fatal("expected timeFromString to return false on 0 len str")
		}
	})
}

func TestBounds(t *testing.T) {
	defer initIntSizes(mathIntSize)

	var c Conv
	chk := func() {
		chkMaxInt, err := c.Int(fmt.Sprintf("%v", math.MaxInt64))
		if err != nil {
			t.Error(err)
		}
		if int64(chkMaxInt) != mathMaxInt {
			t.Errorf("chkMaxInt exp %v; got %v", chkMaxInt, mathMaxInt)
		}

		chkMinInt, err := c.Int(fmt.Sprintf("%v", math.MinInt64))
		if err != nil {
			t.Error(err)
		}
		if int64(chkMinInt) != mathMinInt {
			t.Errorf("chkMaxInt exp %v; got %v", chkMinInt, mathMaxInt)
		}

		chkUint, err := c.Uint(fmt.Sprintf("%v", uint64(math.MaxUint64)))
		if err != nil {
			t.Error(err)
		}
		if uint64(chkUint) != mathMaxUint {
			t.Errorf("chkMaxInt exp %v; got %v", chkMinInt, chkUint)
		}
	}

	initIntSizes(32)
	chk()

	initIntSizes(64)
	chk()
}
