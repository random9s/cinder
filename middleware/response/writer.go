package response

import (
	"net/http"
	"strconv"
)

type responseHandler struct {
	ContentEncoding string
	h               http.Handler
}

//Writer writes data to the responsewriter
func Writer(contentEncoding string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &responseHandler{contentEncoding, h}
	}
}

func (h *responseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := New()

	h.h.ServeHTTP(resp, r)

	enc := encoder(h.ContentEncoding)
	b, err := enc(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for key, val := range resp.header {
		w.Header()[key] = val
	}

	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(b)), 10))
	w.WriteHeader(resp.status)
	w.Write(b)
}
