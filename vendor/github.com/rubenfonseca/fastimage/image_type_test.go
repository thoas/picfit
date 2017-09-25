package fastimage

import "testing"

func TestImageTypeName(t *testing.T) {
	if GIF.String() != "GIF" {
		t.Error("Bad GIF image iname")
	}

	if PNG.String() != "PNG" {
		t.Error("Bad PNG image name")
	}

	if JPEG.String() != "JPEG" {
		t.Error("Bad JPEG image name")
	}
}
