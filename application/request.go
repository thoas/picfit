package application

import (
	"github.com/gorilla/mux"
	"github.com/thoas/gokvstores"
	"github.com/thoas/picfit/engines"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/signature"
	"github.com/thoas/picfit/util"
	"net/http"
	"net/url"
)

type Request struct {
	Request     *http.Request
	Operation   *engines.Operation
	Connection  gokvstores.KVStoreConnection
	Key         string
	URL         *url.URL
	Filepath    string
	Params      map[string]string
	QueryString map[string]string
}

const SIG_PARAM_NAME = "sig"

func NewRequest(req *http.Request, con gokvstores.KVStoreConnection) (*Request, error) {
	req.ParseForm()

	queryString := make(map[string]string)
	params := mux.Vars(req)

	if len(req.Form) > 0 {
		for k, v := range req.Form {
			queryString[k] = v[0]
		}
	}

	for k, v := range params {
		queryString[k] = v
	}

	request := &Request{
		QueryString: queryString,
		Params:      params,
	}

	extracted := map[string]interface{}{}

	for key, extractor := range Extractors {
		result, err := extractor(key, request)

		if err != nil {
			return nil, err
		}

		extracted[key] = result
	}

	sorted := util.SortMapString(queryString)

	delete(sorted, SIG_PARAM_NAME)

	serialized := hash.Serialize(sorted)

	key := hash.Tokey(serialized)

	var u *url.URL
	var path string

	value, ok := extracted["url"]

	if ok && value != nil {
		u = value.(*url.URL)
	}

	value, ok = extracted["path"]

	if ok && value != nil {
		path = value.(string)
	}

	request.Operation = extracted["op"].(*engines.Operation)
	request.Connection = con
	request.Key = key
	request.URL = u
	request.Filepath = path

	return request, nil
}

func (r *Request) IsAuthorized(key string) bool {
	params := url.Values{}
	for k, v := range util.SortMapString(r.QueryString) {
		params.Set(k, v)
	}

	return signature.VerifySign(key, params.Encode())
}
