package fastimage

import "testing"

var buffer = []byte{0x42, 0x24}

func TestReadUint16(t *testing.T) {
	number := readUint16(buffer)
	if number != 16932 {
		t.Error("Error converting bytes to big-endian uint16")
	}
}

func TestReadULint16(t *testing.T) {
	number := readULint16(buffer)
	if number != 9282 {
		t.Error("Error converting bytes to little-endian uint16")
	}
}
