package application

import (
	"fmt"
)

func panicIf(err error) {
	if err != nil {
		App.Logger.Error.Print(err)

		panic(err)
	}
}

func mapInterfaceToMapString(obj map[string]interface{}) map[string]string {
	results := make(map[string]string)

	for k, v := range obj {
		results[k] = fmt.Sprintf("%v", v)
	}

	return results
}
