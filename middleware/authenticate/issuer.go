package authenticate

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

//Auth ...
type Auth string

//Key ...
const Key Auth = "authenticate:user"

//Authenticated confirms if a user was authenticated
type Authenticated struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
	Status  int    `json:"status"`
	Err     string `json:"error"`

	authenticated bool
	identifier    string
	subject       string
	signer        []byte
	audience      string
	issuer        string
}

//NewAuthenticated ...
func NewAuthenticated() *Authenticated {
	return &Authenticated{}
}

//Identifier sets the identifier for the token being issued
func (a *Authenticated) Identifier(id string) {
	a.identifier = id
}

//Subject ...
func (a *Authenticated) Subject(sub string) {
	a.subject = sub
}

//Signer sets the signature that will encrypt tokens
func (a *Authenticated) Signer(key []byte) {
	a.signer = key
}

//Audience tells the server who the jwt is intended for
func (a *Authenticated) Audience(aud string) {
	a.audience = aud
}

//Issuer tells the server who issued the token
func (a *Authenticated) Issuer(iss string) {
	a.issuer = iss
}

//Success satisfies the response interface
func (a *Authenticated) Success(v interface{}) error {
	a.authenticated = true
	return nil
}

//PlainText satisfies the response interface
func (a *Authenticated) PlainText(b []byte) error {
	a.authenticated = true
	return nil
}

//Error ...
func (a *Authenticated) Error(err error, status int) error {
	a.authenticated = false
	a.Err = err.Error()
	a.Status = status
	return nil
}

//Issuer is responsible for issuing JWT after a user has been authenticated
type Issuer struct {
	Exp           int64
	Hijacker      func(string, http.ResponseWriter, *http.Request) error
	SigningMethod jwt.SigningMethod
}

//NewIssuer ...
func NewIssuer(exp int64, hijacker func(string, http.ResponseWriter, *http.Request) error, signingMethod jwt.SigningMethod) *Issuer {
	return &Issuer{
		Exp:           exp,
		Hijacker:      hijacker,
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
	ctx := context.WithValue(r.Context(), Key, NewAuthenticated())
	r = r.WithContext(ctx)

	defer func(req *http.Request) {
		//Check that authorized value was set
		v := r.Context().Value(Key)
		if v == nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("context value not found"))
			return
		}

		//Check that authorized value is type bool
		valid, ok := v.(*Authenticated)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("invalid type assertion"))
			return
		}

		//Check if user was authenticated successfully
		if !valid.authenticated {
			issueError(valid, w)
			return
		}

		//Create JWT and send
		c := NewClaim(h.issuer.Exp, valid)
		token, err := c.NewJWT(h.issuer.SigningMethod, valid.signer)
		if err != nil {
			valid.Err = err.Error()
			issueError(valid, w)
			return
		}

		if h.issuer.Hijacker != nil {
			err = h.issuer.Hijacker(token, w, r)
			if err != nil {
				valid.Err = "issuer hijack func err: " + err.Error()
				issueError(valid, w)
				return
			}
		}

		valid.Access = fmt.Sprintf("Bearer %v", token)
		issueToken(valid, w)
	}(r)

	h.h.ServeHTTP(w, r)
}

func issueToken(a *Authenticated, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Content-Encoding", "gzip")

	gw := gzip.NewWriter(w)
	defer gw.Close()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(gw).Encode(a)
}

func issueError(a *Authenticated, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Content-Encoding", "gzip")

	gw := gzip.NewWriter(w)
	defer gw.Close()

	w.WriteHeader(a.Status)
	json.NewEncoder(gw).Encode(a)
}
