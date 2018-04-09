package response

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"context"
	"net/http"
)

//Resp ...
type Resp string

//Key ...
const Key Resp = "wrappedResponse"

//Response compressions
const (
	GZIP     = "gzip"
	LZW      = "compress"
	FLATE    = "deflate"
	BROTLI   = "br"
	IDENTITY = "identity"
)

type responseHandler struct {
	contentEncoding string
	h               http.Handler
}

//Writer writes data to the responsewriter
func Writer(contentEncoding string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &responseHandler{contentEncoding, h}
	}
}

func (h *responseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := New(w)
	resp.setEncoder(h.contentEncoding)

	ctx := context.WithValue(r.Context(), Key, resp)
	h.h.ServeHTTP(w, r.WithContext(ctx))
}

func (r *Response) setEncoder(ce string) {
	switch ce {
	case GZIP:
		r.encoder = gzipcompress
	case LZW:
		r.encoder = lzwcompress
	case FLATE:
		r.encoder = deflate
	default:
		r.encoder = identity
	}
}

func gzipcompress(r *Response) ([]byte, error) {
	r.w.Header().Set("Content-Encoding", "gzip")

	var buff bytes.Buffer
	writer := gzip.NewWriter(&buff)
	defer writer.Close()

	_, err := writer.Write(r.Bytes())
	return buff.Bytes(), err
}

func lzwcompress(r *Response) ([]byte, error) {
	r.w.Header().Set("Content-Encoding", "compress")

	var buff bytes.Buffer
	writer := lzw.NewWriter(&buff, lzw.MSB, 8)
	defer writer.Close()

	_, err := writer.Write(r.Bytes())
	return buff.Bytes(), err
}

func deflate(r *Response) ([]byte, error) {
	r.w.Header().Set("Content-Encoding", "deflate")

	var buff bytes.Buffer
	writer, err := flate.NewWriter(&buff, -1)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	_, err = writer.Write(r.Bytes())
	return buff.Bytes(), err
}

func identity(r *Response) ([]byte, error) {
	r.w.Header().Set("Content-Encoding", "identity")
	return r.Bytes(), nil
}
