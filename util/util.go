package util

import (
	"fmt"
	"sort"
)

func MapInterfaceToMapString(obj map[string]interface{}) map[string]string {
	results := make(map[string]string)

	for k, v := range obj {
		results[k] = fmt.Sprintf("%v", v)
	}

	return results
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
