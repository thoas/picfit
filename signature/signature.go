package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
)

func Sign(key string, qs string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(qs))

	byteArray := mac.Sum(nil)

	return hex.EncodeToString(byteArray)
}

func AppendSign(key string, qs string) string {
	signature := Sign(key, qs)

	params := url.Values{}
	params.Add("sig", signature)

	return fmt.Sprintf("%s&%s", qs, params.Encode())
}

func VerifySign(key string, qs string) bool {
	r, _ := regexp.Compile("&?sig=[^&]*")

	unsignedQueryString := r.ReplaceAllString(qs, "")

	sign := Sign(key, unsignedQueryString)
	values, _ := url.ParseQuery(qs)

	return values.Get("sig") == sign
}
