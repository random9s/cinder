package response

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
	fmt.Println("check if written")
	if !r.Written() {
		r.Body = v
		r.toJSON()
		r.Status = http.StatusOK
		r.ContentType = "application/json; charset=UTF-8"
		return r.write()
	}

	return errors.New("response has already been written")
}

//WriteXML returns an XML Encdoded server response
func (r *Response) WriteXML(v interface{}) error {
	if !r.Written() {
		r.Status = http.StatusOK
		r.toXML()
		r.Status = http.StatusOK
		r.ContentType = "application/xml; charset=UTF-8"
		return r.write()
	}

	return errors.New("response has already been written")
}

func (r *Response) write() error {
	r.w.Header().Set("Accept-Charset", "utf-8")
	r.w.Header().Set("Content-Type", fmt.Sprintf("%s", r.ContentType))

	b, err := r.encoder(r)
	if err != nil {
		r.w.WriteHeader(http.StatusInternalServerError)
		r.w.Write([]byte(err.Error()))
		return err
	}

	r.w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
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
	var buff = &bytes.Buffer{}

	err := xml.NewEncoder(buff).Encode(r.Body)
	if err != nil {
		r.Status = http.StatusInternalServerError
		r.body = bytes.NewBuffer([]byte(err.Error()))
		return
	}

	r.body = buff
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
