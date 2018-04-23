package response

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
)

//Response ...
type Response struct {
	encoder func(*Response) ([]byte, error)
	body    *bytes.Buffer

	header http.Header
	status int
}

//New wraps responsewriter
func New() *Response {
	return &Response{
		body:   &bytes.Buffer{},
		header: make(http.Header),
	}
}

//Header returns http headers
func (r *Response) Header() http.Header {
	return r.header
}

//WriteHeader ...
func (r *Response) WriteHeader(statusCode int) {
	r.status = statusCode
}

//Status ...
func (r *Response) Status() int {
	return r.status
}

//WriteJSON returns a JSON Encoded server response
func (r *Response) WriteJSON(v interface{}) (int, error) {
	if !r.Written() {
		b, err := json.Marshal(v)
		if err != nil {
			return 0, err
		}

		if r.header.Get("Content-Type") == "" {
			r.header.Add("Content-Type", "application/json; charset=UTF-8")
		}
		return r.Write(b)
	}

	return 0, errors.New("response has already been written")
}

//WriteXML returns an XML Encdoded server response
func (r *Response) WriteXML(v interface{}) (int, error) {
	if !r.Written() {
		b, err := xml.Marshal(v)
		if err != nil {
			return 0, err
		}

		if r.header.Get("Content-Type") == "" {
			r.header.Add("Content-Type", "application/xml; charset=UTF-8")
		}

		return r.Write(b)
	}

	return 0, errors.New("response has already been written")
}

//Error ...
func (r *Response) Error(err error, status int) (int, error) {
	if !r.Written() {
		var e = &struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}

		r.status = status
		b, err := json.Marshal(e)
		if err != nil {
			return 0, err
		}

		return r.Write(b)
	}

	return 0, errors.New("response has already been written")
}

//Write ...
func (r *Response) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

//Len of the response body
func (r *Response) Len() int {
	return r.body.Len()
}

//Bytes returns body to write
func (r *Response) Bytes() []byte {
	return r.body.Bytes()
}

//Written checks if response has been written
func (r *Response) Written() bool {
	return r.Len() > 0
}
