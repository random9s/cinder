package authenticate

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
)

//Token confirms if a user was authenticated
type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`

	identifier string
	subject    string
	audience   string
	issuer     string

	exp        int64
	signMethod jwt.SigningMethod

	signer []byte
}

//NewToken ...
func NewToken(exp int64, method jwt.SigningMethod, signer []byte) *Token {
	return &Token{
		exp:        exp,
		signer:     signer,
		signMethod: method,
	}
}

//Identifier sets the identifier for the token being issued
func (t *Token) Identifier(id string) {
	t.identifier = id
}

//Subject ...
func (t *Token) Subject(sub string) {
	t.subject = sub
}

//Audience tells the server who the jwt is intended for
func (t *Token) Audience(aud string) {
	t.audience = aud
}

//Issuer tells the server who issued the token
func (t *Token) Issuer(iss string) {
	t.issuer = iss
}

//Signer ...
func (t *Token) Signer(signer []byte) {
	t.signer = signer
}

//SignMethod sets the signing method
func (t *Token) SignMethod(signingMethod jwt.SigningMethod) {
	t.signMethod = signingMethod
}

//Expires sets the expire date
func (t *Token) Expires(exp int64) {
	t.exp = exp
}

//GenerateAccess returns access token
func GenerateAccess(t *Token) error {
	if t.signMethod == nil {
		return errors.New("must set signing method")
	}

	if t.signer == nil {
		return errors.New("must set signing key")
	}

	//Create JWT and send
	var err error
	t.Access, err = NewClaim(t).NewJWT(t.signMethod, t.signer)
	return err
}

//GenerateRefresh returns refresh token
func GenerateRefresh(t *Token) error {
	if t.signMethod == nil {
		return errors.New("must set signing method")
	}

	if t.signer == nil {
		return errors.New("must set signing key")
	}

	//Create JWT and send
	var err error
	t.Refresh, err = NewClaim(t).NewJWT(t.signMethod, t.signer)
	return err
}
