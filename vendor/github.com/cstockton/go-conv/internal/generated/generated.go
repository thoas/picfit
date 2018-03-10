package generated

import "time"

var (
	expBoolVal1     = true
	expBoolVal2     = false
	expDurationVal1 = time.Duration(10)       // Duration("10ns")
	expDurationVal2 = time.Duration(20000)    // Duration("20µs")
	expDurationVal3 = time.Duration(30000000) // Duration("30ms")
	expFloat32Val1  = float32(1.2)
	expFloat32Val2  = float32(3.45)
	expFloat32Val3  = float32(6.78)
	expFloat64Val1  = float64(1.2)
	expFloat64Val2  = float64(3.45)
	expFloat64Val3  = float64(6.78)
	expIntVal1      = int(12)
	expIntVal2      = int(34)
	expIntVal3      = int(56)
	expInt16Val1    = int16(12)
	expInt16Val2    = int16(34)
	expInt16Val3    = int16(56)
	expInt32Val1    = int32(12)
	expInt32Val2    = int32(34)
	expInt32Val3    = int32(56)
	expInt64Val1    = int64(12)
	expInt64Val2    = int64(34)
	expInt64Val3    = int64(56)
	expInt8Val1     = int8(12)
	expInt8Val2     = int8(34)
	expInt8Val3     = int8(56)
	expStringVal1   = "k1"
	expStringVal2   = "K2"
	expStringVal3   = "03"
	expUintVal1     = uint(12)
	expTimeVal1     = time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	expTimeVal2     = time.Date(2006, 1, 2, 16, 4, 5, 0, time.UTC)
	expTimeVal3     = time.Date(2006, 1, 2, 17, 4, 5, 0, time.UTC)
	// expTimeVal1   = mustTime("2006-01-02T15:04:05Z")
	// expTimeVal2   = mustTime("2006-01-02T16:04:05Z")
	// expTimeVal3   = mustTime("2006-01-02T17:04:05Z")
	expUintVal2   = uint(34)
	expUintVal3   = uint(56)
	expUint16Val1 = uint16(12)
	expUint16Val2 = uint16(34)
	expUint16Val3 = uint16(56)
	expUint32Val1 = uint32(12)
	expUint32Val2 = uint32(34)
	expUint32Val3 = uint32(56)
	expUint64Val1 = uint64(12)
	expUint64Val2 = uint64(34)
	expUint64Val3 = uint64(56)
	expUint8Val1  = uint8(12)
	expUint8Val2  = uint8(34)
	expUint8Val3  = uint8(56)

	strToNumeric = [42]struct {
		from string
		to   int64
	}{
		{"0", 0},
		{"-0", 0},
		{"1", 1},
		{"-1", -1},
		{"12", 12},
		{"-12", -12},
		{"123", 123},
		{"-123", -123},
		{"1234", 1234},
		{"-1234", -1234},
		{"12345", 12345},
		{"-12345", -12345},
		{"123456", 123456},
		{"-123456", -123456},
		{"1234567", 1234567},
		{"-1234567", -1234567},
		{"12345678", 12345678},
		{"-12345678", -12345678},
		{"123456789", 123456789},
		{"-123456789", -123456789},
		{"1234567890", 1234567890},
		{"-1234567890", -1234567890},
		{"12345678901", 12345678901},
		{"-12345678901", -12345678901},
		{"123456789012", 123456789012},
		{"-123456789012", -123456789012},
		{"1234567890123", 1234567890123},
		{"-1234567890123", -1234567890123},
		{"12345678901234", 12345678901234},
		{"-12345678901234", -12345678901234},
		{"123456789012345", 123456789012345},
		{"-123456789012345", -123456789012345},
		{"1234567890123456", 1234567890123456},
		{"-1234567890123456", -1234567890123456},
		{"12345678901234567", 12345678901234567},
		{"-12345678901234567", -12345678901234567},
		{"123456789012345678", 123456789012345678},
		{"-123456789012345678", -123456789012345678},
		{"1234567890123456789", 1234567890123456789},
		{"-1234567890123456789", -1234567890123456789},
		{"9223372036854775807", 1<<63 - 1},
		{"-9223372036854775808", -1 << 63},
	}
)
