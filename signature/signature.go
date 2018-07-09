package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
)

var signRegex = regexp.MustCompile("&?sig=[^&]*")

// VerifyParameters encodes map parameters with a key and returns if parameters match signature
func VerifyParameters(key string, qs map[string]string) bool {
	params := url.Values{}

	for k, v := range qs {
		params.Set(k, v)
	}

	return VerifySign(key, params.Encode())
}

// Sign encodes query string using a key
func Sign(key string, qs string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(qs))

	byteArray := mac.Sum(nil)

	return hex.EncodeToString(byteArray)
}

// SignRaw encodes raw query string (not sorted) using a key
func SignRaw(key string, queryString string) (string, error) {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return "", err
	}

	return Sign(key, values.Encode()), nil
}

// AppendSign appends the signature to query string
func AppendSign(key string, qs string) string {
	signature := Sign(key, qs)

	params := url.Values{}
	params.Add("sig", signature)

	return fmt.Sprintf("%s&%s", qs, params.Encode())
}

// VerifySign extracts the signature and compare it with query string
func VerifySign(key string, qs string) bool {
	unsignedQueryString := signRegex.ReplaceAllString(qs, "")

	sign := Sign(key, unsignedQueryString)
	values, _ := url.ParseQuery(qs)

	return values.Get("sig") == sign
}
