package storage

import (
	"fmt"

	"github.com/mitchellh/goamz/aws"
)

var regions = map[string]struct{}{
	"ams2": {},
	"ams3": {},
	"blr1": {},
	"fra1": {},
	"lon1": {},
	"nyc1": {},
	"nyc2": {},
	"nyc3": {},
	"sfo1": {},
	"sfo2": {},
	"sgp1": {},
	"tor1": {},
}

var GetDOs3Region = func(region string) (aws.Region, bool) {
	if _, ok := regions[region]; !ok {
		return aws.Region{}, false
	}

	return aws.Region{
		Name:       region,
		S3Endpoint: fmt.Sprintf("https://%s.digitaloceanspaces.com", region),
	}, true
}
