package authenticate

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

//Auth ...
type Auth string

//Key ...
const Key Auth = "authenticate:user"

//Issuer is responsible for issuing JWT after a user has been authenticated
type Issuer struct {
	Exp           int64
	SigningMethod jwt.SigningMethod
}

//NewIssuer ...
func NewIssuer(exp int64, signingMethod jwt.SigningMethod) *Issuer {
	return &Issuer{
		Exp:           exp,
		SigningMethod: signingMethod,
	}
}

type issuerHandler struct {
	h http.Handler

	issuer *Issuer
}

//IssueJWT creates an issuer that distributes JWT to authenticated users
func IssueJWT(issuer *Issuer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &issuerHandler{
			h:      h,
			issuer: issuer,
		}
	}
}

func (h *issuerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), Key, NewToken(h.issuer.Exp, h.issuer.SigningMethod))
	h.h.ServeHTTP(w, r.WithContext(ctx))
}
