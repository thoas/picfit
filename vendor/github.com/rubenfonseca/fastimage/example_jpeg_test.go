package fastimage_test

import (
	"fmt"
	"os"

	"github.com/rubenfonseca/fastimage"
)

// This example shows basic usage of the package: just pass an url to the
// detector, and analyze the results.
func Example_remoteBigJPEG() {
	url := "http://upload.wikimedia.org/wikipedia/commons/9/9a/SKA_dishes_big.jpg"

	imagetype, size, err := fastimage.DetectImageType(url)
	if err != nil {
		// Something went wrong, http failed? not an image?
		panic(err)
	}

	fmt.Printf("Image size: %v\n", size)

	switch imagetype {
	case fastimage.JPEG:
		fmt.Println("JPEG")
	case fastimage.PNG:
		fmt.Println("PNG")
	case fastimage.GIF:
		fmt.Println("GIF")
	}
	// Output: Image size: &{5000 2813}
	// JPEG
}

func Example_localBigJPEG() {
	f, err := os.Open("example.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	imagetype, size, err := fastimage.DetectImageTypeFromReader(f)
	if err != nil {
		// Something went wrong, not an image?
		panic(err)
	}

	fmt.Printf("Image size: %v\n", size)

	switch imagetype {
	case fastimage.JPEG:
		fmt.Printf("JPEG")
	case fastimage.PNG:
		fmt.Printf("PNG")
	case fastimage.GIF:
		fmt.Printf("GIF")
	}
	// Output: Image size: &{320 240}
	// GIF
}
