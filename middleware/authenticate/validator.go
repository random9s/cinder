package authenticate

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

//VerifyFunc ...
type VerifyFunc func(*jwt.Token) error

//SecretFunc ...
type SecretFunc func(jwt.Claims) (interface{}, error)

//Validator expresses how the JWT should be parsed and validated
type Validator struct {
	Claims        jwt.Claims        //Custom claim to decode into
	SigningMethod jwt.SigningMethod //SigningMethod to decode jwt

	verify VerifyFunc // User input validation for claims
	secret SecretFunc // user input to retrieve secret key

	token *jwt.Token //Populated when ValidateToken is called
}

//NewValidator returns initialized validator
func NewValidator(claims jwt.Claims, signingMethod jwt.SigningMethod, verify VerifyFunc, secret SecretFunc) *Validator {
	return &Validator{
		Claims:        claims,
		SigningMethod: signingMethod,
		verify:        verify,
		secret:        secret,
	}
}

//ValidateToken retrieves and parses token to check for validity
func (v *Validator) ValidateToken(r *http.Request) error {
	//Retrieve token from header
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		return errors.New("empty authorization header")
	}

	//Remove Bearer from token
	token := strings.TrimPrefix(bearer, "Bearer ")

	//Set key func
	var keyFunc func(*jwt.Token) (interface{}, error)

	switch v.SigningMethod.(type) {
	case *jwt.SigningMethodHMAC:
		keyFunc = v.methodHMACKeyFunc()
	case *jwt.SigningMethodECDSA:
		keyFunc = v.methodECDSAKeyFunc()
	case *jwt.SigningMethodRSA:
		keyFunc = v.methodRSAKeyFunc()
	default:
		return errors.New("signing method is required")
	}

	//return parsed token
	parsedToken, err := jwt.ParseWithClaims(token, v.Claims, keyFunc)
	if err != nil {
		return err
	}

	if !parsedToken.Valid {
		return errors.New("invalid token provided")
	}

	v.token = parsedToken
	return nil
}

//ValidateClaims checks claims!
func (v *Validator) ValidateClaims() error {
	return v.verify(v.token)
}

// ~~~~~~~~~~~~~~ SIGNING METHOD KEY FUNCS ~~~~~~~~~~~~~~ //
func (v *Validator) methodHMACKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return v.secret(token.Claims)
	}
}

func (v *Validator) methodECDSAKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return v.secret(token.Claims)
	}
}

func (v *Validator) methodRSAKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		return v.secret(token.Claims)
	}
}
