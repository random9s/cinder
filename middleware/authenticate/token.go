package authenticate

import jwt "github.com/dgrijalva/jwt-go"

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
func NewToken(exp int64, method jwt.SigningMethod) *Token {
	return &Token{
		exp:        exp,
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

//GenerateAccess returns access token
func GenerateAccess(t *Token) error {
	//Create JWT and send
	var err error
	t.Access, err = NewClaim(t).NewJWT(t.signMethod, t.signer)
	return err
}

//GenerateRefresh returns refresh token
func GenerateRefresh(t *Token) error {
	//Create JWT and send
	var err error
	t.Refresh, err = NewClaim(t).NewJWT(t.signMethod, t.signer)
	return err
}
