package authenticate

import (
	"net/http"
)

type authenticationHandler struct {
	Secure bool

	h    http.Handler
	fail http.Handler

	validator *Validator
}

//JWT checks for authenticity of a JWT
func JWT(validator *Validator, secure bool, fail http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &authenticationHandler{
			h:         h,
			fail:      fail,
			Secure:    secure,
			validator: validator,
		}
	}
}

func (h *authenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Secure {
		//Parse token from request
		err := h.validator.ValidateToken(r)
		if err != nil {
			h.fail.ServeHTTP(w, r)
			return
		}

		err = h.validator.ValidateClaims()
		if err != nil {
			h.fail.ServeHTTP(w, r)
			return
		}
	}

	h.h.ServeHTTP(w, r)
}
