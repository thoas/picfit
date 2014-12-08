package hash

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strings"
)

func Tokey(args ...string) string {
	hasher := md5.New()
	hasher.Write([]byte(strings.Join(args, "||")))

	return hex.EncodeToString(hasher.Sum(nil))
}

func Serialize(obj interface{}) string {
	result, _ := json.Marshal(obj)

	return string(result)
}

func Shard(str string, width int, depth int, restOnly bool) []string {
	var results []string

	for i := 0; i < depth; i++ {
		results = append(results, str[(width*i):(width*(i+1))])
	}

	if restOnly {
		results = append(results, str[(width*depth):])
	} else {
		results = append(results, str)
	}

	return results
}
