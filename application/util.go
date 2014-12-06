package application

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

func panicIf(err error) {
	if err != nil {
		App.Logger.Error.Print(err)

		panic(err)
	}
}

func tokey(args ...string) string {
	hasher := md5.New()
	hasher.Write([]byte(strings.Join(args, "||")))

	return hex.EncodeToString(hasher.Sum(nil))
}

func serialize(obj interface{}) string {
	result, _ := json.Marshal(obj)

	return string(result)
}

func shard(str string, width int, depth int, restOnly bool) []string {
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

func mapInterfaceToMapString(obj map[string]interface{}) map[string]string {
	results := make(map[string]string)

	for k, v := range obj {
		results[k] = fmt.Sprintf("%v", v)
	}

	return results
}
