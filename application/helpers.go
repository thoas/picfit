package application

import (
	"github.com/thoas/muxer"
)

func keyFromRequest(req *muxer.Request) string {
	qs := serialize(req.QueryString)

	var key string

	if filename, ok := req.Params["filename"]; !ok {
		key = tokey(req.Params["method"], qs)
	} else {
		key = tokey(req.Params["method"], filename, qs)
	}

	return key
}
