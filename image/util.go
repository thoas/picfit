package image

func matchInArray(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func mapKeys(mapping map[string][]string) []string {
	mk := make([]string, len(mapping))

	i := 0

	for k, _ := range mapping {
		mk[i] = k
		i++
	}

	return mk
}
