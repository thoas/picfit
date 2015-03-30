package application

import (
	"github.com/thoas/gokvstores"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/engines"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/signature"
	"github.com/thoas/picfit/util"
	"net/http"
	"net/url"
)

type Request struct {
	*muxer.Request
	Operation  *engines.Operation
	Connection gokvstores.KVStoreConnection
	Key        string
	URL        *url.URL
	Filepath   string
}

const SIG_PARAM_NAME = "sig"

func NewRequest(req *http.Request, con gokvstores.KVStoreConnection) (*Request, error) {
	request := muxer.NewRequest(req)

	for k, v := range request.Params {
		request.QueryString[k] = v
	}

	extracted := map[string]interface{}{}

	for key, extractor := range Extractors {
		result, err := extractor(key, request)

		if err != nil {
			return nil, err
		}

		extracted[key] = result
	}

	sorted := util.SortMapString(request.QueryString)

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

	return &Request{
		request,
		extracted["op"].(*engines.Operation),
		con,
		key,
		u,
		path,
	}, nil
}

func (r *Request) IsAuthorized(key string) bool {
	params := url.Values{}
	for k, v := range util.SortMapString(r.QueryString) {
		params.Set(k, v)
	}

	return signature.VerifySign(key, params.Encode())
}
