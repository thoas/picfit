package application

import (
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/hash"
)

func keyFromRequest(req *muxer.Request) string {
	qs := hash.Serialize(req.QueryString)

	var key string

	if filename, ok := req.Params["filename"]; !ok {
		key = hash.Tokey(req.Params["method"], qs)
	} else {
		key = hash.Tokey(req.Params["method"], filename, qs)
	}

	return key
}
