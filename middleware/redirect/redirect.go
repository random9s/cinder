package redirect

import (
	"net/http"
)

type redirectHandler struct {
	h       http.Handler
	enabled bool
}

//OnHTTPS checks for redirection.  Enable is an easy switch for debug/development environments
func OnHTTPS(enable bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &redirectHandler{
			h:       h,
			enabled: enable,
		}
	}
}

func (h *redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.enabled {
		if r.Proto == "HTTP/1.1" {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		}
	}

	h.h.ServeHTTP(w, r)
}
