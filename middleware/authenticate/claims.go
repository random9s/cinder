package authenticate

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//CustomClaims ...
type CustomClaims struct {
	jwt.StandardClaims
}

//NewClaim returns new claims
func NewClaim(t *Token) *CustomClaims {
	var iat = time.Now().Unix()
	var expTime = iat + t.exp

	return &CustomClaims{
		jwt.StandardClaims{
			NotBefore: iat,
			IssuedAt:  iat,
			ExpiresAt: expTime,
			Id:        t.identifier,
			Subject:   t.subject,
			Issuer:    t.issuer,
			Audience:  t.audience,
		},
	}
}

//NewJWT returns a new jwt
func (c *CustomClaims) NewJWT(signingMethod jwt.SigningMethod, secret []byte) (string, error) {
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}

//Verify claims
func (c *CustomClaims) Verify() error {
	vErr := new(jwt.ValidationError)
	now := jwt.TimeFunc().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if c.VerifyExpiresAt(now, false) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if c.VerifyIssuedAt(now, false) == false {
		vErr.Inner = fmt.Errorf("Token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if c.VerifyNotBefore(now, false) == false {
		vErr.Inner = fmt.Errorf("token is not valid yet")
		vErr.Errors |= jwt.ValidationErrorNotValidYet
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}
