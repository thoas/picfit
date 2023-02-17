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

func SortMapString(obj map[string]interface{}) map[string]interface{} {
	mk := make([]string, len(obj))

	i := 0
	for k, _ := range obj {
		mk[i] = k
		i++
	}

	sort.Strings(mk)

	results := make(map[string]interface{})

	for _, index := range mk {
		results[index] = obj[index]
	}

	return results
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
