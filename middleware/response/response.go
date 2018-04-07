package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//Response ...
type Response struct {
	Body        interface{}
	ContentType string
	Status      int

	encoder func(*Response) ([]byte, error)
	body    *bytes.Buffer

	w       http.ResponseWriter
	written bool
}

//New wraps responsewriter
func New(w http.ResponseWriter) *Response {
	return &Response{
		Body:    nil,
		body:    &bytes.Buffer{},
		w:       w,
		written: false,
	}
}

//WriteJSON returns a JSON Encoded server response
func (r *Response) WriteJSON(v interface{}) error {
	if !r.Written() {
		r.Body = v
		r.toJSON()
		r.ContentType = "application/json; charset=UTF-8"
		return r.write()
	}

	return errors.New("response has already been written")
}

//WriteXML returns an XML Encdoded server response
func (r *Response) WriteXML(v interface{}) error {
	//TODO
	return errors.New("xml is not yet supported")
}

func (r *Response) write() error {
	r.w.Header().Add("Accept-Charset", "utf-8")
	r.w.Header().Add("Content-Type", fmt.Sprintf("%s", r.ContentType))

	b, err := r.encoder(r)
	if err != nil {
		return err
	}

	r.w.Header().Add("Content-Length", fmt.Sprintf("%s", len(b)))
	r.w.WriteHeader(r.Status)
	r.w.Write(b)
	r.written = true
	return nil
}

//Error ...
func (r *Response) Error(err error, status int) error {
	if !r.Written() {
		var e = &struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}

		r.Body = e
		r.Status = status
		r.toJSON()
		r.ContentType = "application/json; charset=UTF-8"

		return r.write()
	}

	return errors.New("response has already been written")
}

func (r *Response) toJSON() {
	var buff = &bytes.Buffer{}

	err := json.NewEncoder(buff).Encode(r.Body)
	if err != nil {
		r.Status = http.StatusInternalServerError
		r.body = bytes.NewBuffer([]byte(err.Error()))
		return
	}

	r.body = buff
}

func (r *Response) toXML() {
	//TODO
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
	return r.written
}
