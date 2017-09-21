# fastimage

[![GoDoc](https://godoc.org/github.com/rubenfonseca/fastimage?status.png)](https://godoc.org/github.com/rubenfonseca/fastimage) [![Build Status](https://travis-ci.org/rubenfonseca/fastimage.svg?branch=master)](http://travis-ci.org/rubenfonseca/fastimage) [![wercker status](https://app.wercker.com/status/99926bf6ba7846259d5f65de5c5b4ed1/m "wercker status")](https://app.wercker.com/project/bykey/99926bf6ba7846259d5f65de5c5b4ed1)

by Ruben Fonseca (@[rubenfonseca](http://twitter.com/rubenfonseca))

Golang implementation of [fastimage](https://pypi.python.org/pypi/fastimage/0.2.1).
Finds the type and/or size of an image given its uri by fetching as little as needed.

## How?

fastimage parses the image data as it is downloaded. As soon as it finds out
the size and type of the image, it stops the download.

## Install

    $ go get github.com/rubenfonseca/fastimage

## Usage

For instance, this is a big 10MB JPEG image on wikipedia:


	url := "http://upload.wikimedia.org/wikipedia/commons/9/9a/SKA_dishes_big.jpg"

    fastimage.Debug()
	imagetype, size, err := fastimage.DetectImageType(url)
	if err != nil {
		// Something went wrong, http failed? not an image?
		panic(err)
	}

	switch imagetype {
	case fastimage.JPEG:
		log.Printf("JPEG")
	case fastimage.PNG:
		log.Printf("PNG")
	case fastimage.GIF:
		log.Printf("GIF")
	}

	log.Printf("Image type: %s", imagetype.String())
	log.Printf("Image size: %v", size)

At the end, you can read something like this:

    Closed after reading just 17863 bytes out of 10001439 bytes

If you want to set request timeout for url:

    // the second argument is request timeout (milliseconds).
    // FYI, DetectImageType() uses default timeout 5000ms.
    imagetype, size, err := fastimage.DetectImageTypeWithTimeout(url, 2000)

## Supported file types

| File type | Can detect type? | Can detect size? |
|-----------|:----------------:|:----------------:|
| PNG       | Yes              | Yes              |
| JPEG      | Yes              | Yes              |
| GIF       | Yes              | Yes              |
| BMP       | Yes              | No               |
| TIFF      | Yes              | No               |


# Project details

### License

fastimage is under MIT license. See the [LICENSE][license] file for details.

[license]: https://github.com/rubenfonseca/fastimage/blob/master/LICENSE
