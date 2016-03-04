package util

import (
	"fmt"
	"sort"
	"os"
	"strings"
)

const ENV_VAR_PREFIX = "$"

func MapInterfaceToMapString(obj map[string]interface{}) map[string]string {
	results := make(map[string]string)

	for k, v := range obj {
		results[k] = fmt.Sprintf("%v", v)
	}

	return results
}

func parseConfigValue(value string) string {
	if strings.HasPrefix(value, ENV_VAR_PREFIX) {
		envVal := os.Getenv(value[1:])
		if envVal == "" {
			fmt.Printf("Warning: No value set for env var: %s \n", value)
		}
		return envVal
	}
	return value
}

func SortMapString(obj map[string]string) map[string]string {
	mk := make([]string, len(obj))

	i := 0
	for k, _ := range obj {
		mk[i] = k
		i++
	}

	sort.Strings(mk)

	results := make(map[string]string)

	for _, index := range mk {
		results[index] = obj[index]
	}

	return results
}
