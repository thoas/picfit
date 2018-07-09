package application

import (
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/image"
)

type Parameters struct {
	Output     *image.ImageFile
	Operations []engine.EngineOperation
}

func NewParameters(input *image.ImageFile, qs map[string]interface{}) (*Parameters, error) {

	return &Parameters{}, nil
}
