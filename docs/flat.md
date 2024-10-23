# Flat

Flat is a method implemented to the engine goimage in order to draw
images on a background image.

This method can be used only with the multiple operation parameter `op`,
in the URL.

 ## Parameters:

* `path`: the foreground image, can be multiple.
* `pos`: the foreground destination as a rectangle
* `color`: the foreground color in Hex (without `#`), default is transparent.


## Usage

The background is defined by the image transformed by the previous
operation. In order to draw an image on the background a position must
be given in the sub-parameter `pos` and the image path in the
sub-parameter `path`.

example:
```
/display?path=path/to/background.png&op=resize&w=100&h=100
    &op=op:flat+path:path/to/foreground.png+pos:60.10.80.30
```


The value of `pos` must be the coordinates of the rectangle defined
according to the [go image package](https://blog.golang.org/go-image-package):
an axis-aligned rectangle on the integer grid, defined by its top-left and
bottom-right Point.

![Rectangle position](https://github.com/thoas/picfit/blob/main/docs/picfit-dst-position.png)

The foreground image is resized in order to fit in the given rectangle
and centered inside.

If several images are given in the same flat operation with the
subparameters path. The rectangle is cut in equal parts, **horizontally** if
the rectangle width `Dx` is superior to its height `Dy` and
**vertically** if it is not the case. Each images are then resized in
order to fit in each parts and centered inside. The order follow the
given order of `path` parameters in the URL.

![Flat multiple images](https://github.com/thoas/picfit/blob/main/docs/picfit-flat.png)

