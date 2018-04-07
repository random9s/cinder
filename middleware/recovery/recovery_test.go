package recovery

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var panicHandler = func(err interface{}, stacktrace []string, r *http.Request) {}

func TestPanic(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description  string
		url          string
		expectedCode int
		handler      http.Handler
	}{
		{
			description:  "no panic, normal",
			url:          "/",
			expectedCode: http.StatusOK,
			handler:      testHandler,
		}, {
			description:  "panic recovered, log error",
			url:          "/",
			expectedCode: http.StatusInternalServerError,
			handler:      testPanicHandler,
		},
	}

	for _, test := range tests {
		ts := httptest.NewServer(Panic(panicHandler)(test.handler))
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

var testPanicHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	panic("test entered test handler, this should not happen")
})

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
