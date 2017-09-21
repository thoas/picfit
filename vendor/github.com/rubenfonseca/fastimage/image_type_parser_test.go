package fastimage

import "testing"

func TestRegisterNilImageType(t *testing.T) {
	f := func() {
		register(nil)
	}
	if funcDidPanic, _ := didPanic(f); !funcDidPanic {
		t.Error("Registering a nil parser should have paniced")
	}
}

func TestRegisterRepeatedType(t *testing.T) {
	f := func() {
		imageType := imageGIF{}
		register(imageType)
	}

	if funcDidPanic, _ := didPanic(f); !funcDidPanic {
		t.Error("Registering a repeated parser should have paniced")
	}
}

func didPanic(f func()) (bool, interface{}) {
	didPanic := false
	var message interface{}
	func() {
		defer func() {
			if message = recover(); message != nil {
				didPanic = true
			}
		}()

		f()
	}()

	return didPanic, message
}
