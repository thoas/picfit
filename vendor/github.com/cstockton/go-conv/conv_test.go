package conv

import (
	"testing"

	"github.com/cstockton/go-conv/internal/testconv"
)

func TestConv(t *testing.T) {
	testconv.RunBoolTests(t, Bool)
	testconv.RunDurationTests(t, Duration)
	testconv.RunFloat32Tests(t, Float32)
	testconv.RunFloat64Tests(t, Float64)
	testconv.RunInferTests(t, Infer)
	testconv.RunIntTests(t, Int)
	testconv.RunInt8Tests(t, Int8)
	testconv.RunInt16Tests(t, Int16)
	testconv.RunInt32Tests(t, Int32)
	testconv.RunInt64Tests(t, Int64)
	testconv.RunStringTests(t, String)
	testconv.RunTimeTests(t, Time)
	testconv.RunUintTests(t, Uint)
	testconv.RunUint8Tests(t, Uint8)
	testconv.RunUint16Tests(t, Uint16)
	testconv.RunUint32Tests(t, Uint32)
	testconv.RunUint64Tests(t, Uint64)

	// Generates README.md based on testdata/README.md.tpl using examples from
	// the example_test.go file.
	testconv.RunReadmeTest(t, `example_test.go`, `conv_test.go`)
}
