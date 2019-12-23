package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var allowedRegions = []string{
	"nyc1",
	"nyc2",
	"nyc3",
	"ams2",
	"ams3",
	"sfo1",
	"sfo2",
	"sgp1",
	"lon1",
	"fra1",
	"tor1",
	"blr1",
}

func Test_GetDOs3Region_AllowedRegions(t *testing.T) {
	for _, region := range allowedRegions {
		t.Run(region, func(st *testing.T) {
			awsRegion, exist := GetDOs3Region(region)

			assert.True(st, exist)
			assert.Equal(st, region, awsRegion.Name)
			assert.Equal(st, fmt.Sprintf("https://%s.digitaloceanspaces.com", region), awsRegion.S3Endpoint)
		})
	}
}

func Test_GetDOs3Region_WrongRegion(t *testing.T) {
	awsRegion, exist := GetDOs3Region("fake1")

	assert.False(t, exist)
	assert.Empty(t, awsRegion.Name)
	assert.Empty(t, awsRegion.S3Endpoint)
}
