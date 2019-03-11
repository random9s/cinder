package response

import (
	"net/http"
	"strconv"
	"strings"
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

	var ae = r.Header.Get("Accept-Encoding")
	var b = resp.Bytes()
	var err error

	w.Header().Set("Content-Type", http.DetectContentType(b))
	if strings.Contains(ae, h.ContentEncoding) {
		enc := encoder(h.ContentEncoding)
		b, err = enc(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Encoding", "identity")
	}

	for key, val := range resp.header {
		w.Header()[key] = val
	}

	w.Header().Add("Vary", "Accept-Encoding")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(b)), 10))
	w.Header().Set("X-Server-Status", strconv.FormatInt(int64(resp.status), 10))

	w.WriteHeader(resp.status)
	w.Write(b)
}
