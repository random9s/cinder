package authenticate

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

const secret = "239-f3j29-ffwfoiewfjoiwefijwefiowjf0-9293j2f32-9f32-9f"
const exp = 30

var method *jwt.SigningMethodHMAC

func TestIssueJWT(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description  string
		url          string
		expectedCode int
		handler      http.Handler
	}{
		{
			description:  "should return jwt",
			url:          "/",
			expectedCode: http.StatusOK,
			handler:      newTestHandler(true),
		},
	}

	for _, test := range tests {
		issuer := NewIssuer(secret, exp, method)

		ts := httptest.NewServer(IssueJWT(issuer)(test.handler))
		defer ts.Close()

		var u bytes.Buffer
		u.WriteString(string(ts.URL))
		res, err := http.Get(u.String())
		assert.NoError(err)
		if res != nil {
			defer res.Body.Close()
		}

		assert.Equal(test.expectedCode, res.StatusCode, test.description)
	}
}

func newTestHandler(authenticated bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		User(r, authenticated)
	})
}
