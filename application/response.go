package application

import (
	"mime"
	"net/http"
	"strings"
)

type Response struct {
	http.ResponseWriter
}

func NewResponse(res http.ResponseWriter) Response {
	return Response{res}
}

// WriteString writes string data into the response object.
func (r *Response) WriteString(content string) {
	r.ResponseWriter.Write([]byte(content))
}

// Abort is a helper method that sends an HTTP header and an optional
// body. It is useful for returning 4xx or 5xx errors.
// Once it has been called, any return value from the handler will
// not be written to the response.
func (r *Response) Abort(status int, body string) {
	r.ResponseWriter.WriteHeader(status)
	r.WriteString(body)
}

// Redirect is a helper method for 3xx redirects.
func (r *Response) Redirect(status int, url_ string) {
	r.ResponseWriter.Header().Set("Location", url_)
	r.Abort(status, "Redirecting to: "+url_)
}

// PermanentRedirect is a helper method for 301 redirect
func (r *Response) PermanentRedirect(url string) {
	r.Redirect(301, url)
}

// Notmodified writes a 304 HTTP response
func (r *Response) NotModified() {
	r.ResponseWriter.WriteHeader(304)
}

// NotFound writes a 404 HTTP response
func (r *Response) NotFound(message string) {
	r.Abort(404, message)
}

// NotFound writes a 200 HTTP response
func (r *Response) Ok(message string) {
	r.Abort(200, message)
}

//Unauthorized writes a 401 HTTP response
func (r *Response) Unauthorized() {
	r.ResponseWriter.WriteHeader(401)
}

//Forbidden writes a 403 HTTP response
func (r *Response) Forbidden() {
	r.ResponseWriter.WriteHeader(403)
}

//NotAllowed writes a 405 HTTP response
func (r *Response) NotAllowed() {
	r.ResponseWriter.WriteHeader(405)
}

//BadRequest writes a 400 HTTP response
func (r *Response) BadRequest() {
	r.ResponseWriter.WriteHeader(400)
}

// ContentType sets the Content-Type header for an HTTP response.
// For example, ctx.ContentType("json") sets the content-type to "application/json"
// If the supplied value contains a slash (/) it is set as the Content-Type
// verbatim. The return value is the content type as it was
// set, or an empty string if none was found.
func (r *Response) ContentType(val string) string {
	var ctype string
	if strings.ContainsRune(val, '/') {
		ctype = val
	} else {
		if !strings.HasPrefix(val, ".") {
			val = "." + val
		}
		ctype = mime.TypeByExtension(val)
	}
	if ctype != "" {
		r.Header().Set("Content-Type", ctype)
	}
	return ctype
}

// SetHeader sets a response header. If `unique` is true, the current value
// of that header will be overwritten . If false, it will be appended.
func (r *Response) SetHeader(hdr string, val string, unique bool) {
	if unique {
		r.Header().Set(hdr, val)
	} else {
		r.Header().Add(hdr, val)
	}
}

// SetHeaders sets response headers. If `unique` is true, the current value
// of that header will be overwritten . If false, it will be appended.
func (r *Response) SetHeaders(headers map[string]string, unique bool) {
	for k, v := range headers {
		r.SetHeader(k, v, unique)
	}
}
