/*
Package fastimage allows you to find the type and/or size of a remote image by
downloading as little as possible.

Why?
Sometimes you need to know the size of a remote image before downloading it.

How?
fastimage parses the image data as it is downloaded. As soon as it finds out
the size and type of the image, it stops the download.
*/
package fastimage
