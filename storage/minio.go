package storage

import (
	"fmt"

	"github.com/mitchellh/goamz/aws"
)

var GetMINIOs3Region = func(region string, endpoint string) (aws.Region, bool) {

	return aws.Region{
		Name:       region,
		S3Endpoint: fmt.Sprintf(endpoint, region),
	}, true
}
