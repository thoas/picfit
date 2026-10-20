package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	imagefile "github.com/thoas/picfit/image"
)

func TestGoImageAvif(t *testing.T) {
	goimage := &GoImage{}

	avif := run(t, goimage, resizeFunc, readFixture(t, "schwarzy.jpg"),
		Options{Width: 100, Upscale: true, Format: imagefile.AVIF, Quality: 60})
	require.GreaterOrEqual(t, len(avif), 12)
	assert.Equal(t, "ftypavif", string(avif[4:12]))
	cfg, format := decodeConfig(t, avif)
	assert.Equal(t, "avif", format)
	assert.Equal(t, 100, cfg.Width)

	// an avif source is decoded as well
	out := run(t, goimage, resizeFunc, avif, Options{Width: 50, Upscale: true, Format: imagefile.PNG})
	cfg, format = decodeConfig(t, out)
	assert.Equal(t, "png", format)
	assert.Equal(t, 50, cfg.Width)
}
