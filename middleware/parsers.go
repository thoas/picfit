package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/thoas/picfit/constants"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/util"
)

var parametersReg = regexp.MustCompile(`(?:(?P<sig>\w+)/)?(?P<op>\w+)/(?:(?P<w>\d+))?x(?:(?P<h>\d+))?/(?P<path>[\w\-/.]+)`)

// ParametersParser matches parameters to query string
func ParametersParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		result := c.Param("parameters")

		if result != "" {
			match := parametersReg.FindStringSubmatch(result)
			if match != nil {
				parameters := make(map[string]interface{})

				results := parametersReg.SubexpNames()

				for i, name := range results {
					if i != 0 && match[i] != "" {
						parameters[name] = match[i]
					}
				}

				c.Set("parameters", parameters)
			}
		} else {
			if c.Query("url") == "" && c.Query("path") == "" {
				c.String(http.StatusBadRequest, "Request should contains parameters or query string")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// KeyParser injects an unique key from query parameters
func KeyParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryString := make(map[string]interface{})

		params, exists := c.Get("parameters")
		if exists {
			queryString = params.(map[string]interface{})
		}

		sorted := setParamsFromURLValues(queryString, c.Request.URL.Query())
		delete(sorted, constants.SigParamName)
		delete(sorted, constants.ForceParamName)

		if len(sorted) != 0 {
			serialized := hash.Serialize(sorted)

			key := hash.Tokey(serialized)

			c.Set("key", key)
			c.Set("parameters", queryString)
		}

		c.Next()
	}
}

func setParamsFromURLValues(params map[string]interface{}, values url.Values) map[string]interface{} {
	for k, v := range values {
		if k != constants.OperationParamName {
			params[k] = v[0]
			continue
		}

		var operations []string
		op, ok := params[k].(string)
		if ok {
			operations = append(operations, op)
		}
		operations = append(operations, v...)

		if len(operations) > 1 {
			params[k] = operations
		} else if len(operations) == 1 {
			params[k] = operations[0]
		}
	}

	return util.SortMapString(params)
}

// URLParser extracts the url query string and add a url.URL to the context
func URLParser(mimetypeDetectorType string) gin.HandlerFunc {
	mimetypeDetector := image.GetMimetypeDetector(mimetypeDetectorType)

	return func(c *gin.Context) {
		value := c.Query("url")

		if value != "" {
			url, err := url.Parse(value)

			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("URL %s is not valid", value))
				c.Abort()
				return
			}

			mimetype, _ := mimetypeDetector(url)

			_, ok := image.Extensions[mimetype]

			if !ok {
				c.String(http.StatusBadRequest, fmt.Sprintf("Mimetype %s is not supported", mimetype))
				c.Abort()
				return
			}

			c.Set("url", url)
			c.Set("mimetype", mimetype)
		}

		c.Next()
	}
}

// OperationParser extracts the operation and add it to the context
func OperationParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		parameters := c.MustGet("parameters").(map[string]interface{})

		operation, ok := parameters[constants.OperationParamName].(string)
		if ok && operation != "" {
			if _, k := engine.Operations[operation]; !k {
				c.String(http.StatusBadRequest, fmt.Sprintf("Invalid method %s or invalid parameters", operation))
				c.Abort()
				return
			}
			c.Set(constants.OperationParamName, operation)
			c.Next()
			return
		}

		operations, ok := parameters[constants.OperationParamName].([]string)
		if !ok || len(operations) == 0 {
			c.String(http.StatusBadRequest, fmt.Sprintf("`%s` parameter or query string cannot be empty", constants.OperationParamName))
			c.Abort()
			return
		}

		for i := range operations {
			_, ok := engine.Operations[operations[i]]
			if !ok {
				params := make(map[string]string)
				for _, p := range strings.Split(operations[i], " ") {
					l := strings.Split(p, ":")
					if len(l) > 1 {
						params[l[0]] = l[1]
					}
				}

				v, exists := params[constants.OperationParamName]
				if !exists {
					c.String(http.StatusBadRequest, fmt.Sprintf("`%s` parameter or query string cannot be empty", constants.OperationParamName))
					c.Abort()
					return
				} else if _, ok := engine.Operations[v]; !ok {
					c.String(http.StatusBadRequest, fmt.Sprintf("Invalid method %s or invalid parameters", operations[i]))
					c.Abort()
					return
				}
			}
		}

		c.Set("op", operations)
		c.Next()
	}
}
