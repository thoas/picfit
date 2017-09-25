package fastimage

import (
	"io/ioutil"
	"log"
	"os"
)

var logger = log.New(ioutil.Discard, "fastimage", log.LstdFlags)

// Debug enables debug logging of the operations done by the library.
// If called, lots of information will be print to stderr.
func Debug() {
	logger = log.New(os.Stderr, "fastimage", log.LstdFlags)
}
