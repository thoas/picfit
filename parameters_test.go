package picfit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thoas/picfit/tests"
)

func TestEngineOperationFromQuery(t *testing.T) {
	op := "op:resize w:123 h:321 upscale:true pos:top q:99"
	processor := tests.NewDummyProcessor()
	operation, err := processor.NewEngineOperationFromQuery(op)
	assert.Nil(t, err)

	assert.Equal(t, operation.Operation.String(), "resize")
	assert.Equal(t, operation.Options.Height, 321)
	assert.Equal(t, operation.Options.Width, 123)
	assert.Equal(t, operation.Options.Position, "top")
	assert.Equal(t, operation.Options.Quality, 99)
	assert.True(t, operation.Options.Upscale)
}
