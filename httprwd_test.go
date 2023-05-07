package httprwd

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SAMPLE    = "Hello world"
	UNCHANGED = "Response written without modification of the body"
)

var (
	delegate ResponseWriterDelegate
	res      *http.Response
	body     []byte
	ts       = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The delegate is meant to be used exclusively inside middleware (is
		// not intended to be thread safe). Here we reset before each test to
		// simulate instantiation with each request, and do not use -race when
		// running tests. If we wanted to do concurrent testing, we'd need to
		// ensure the delegate is fully local to each request.
		delegate.Code = 0

		// Allow an empty code being passed, so we can test the unset case
		code, err := strconv.Atoi(r.URL.Path[1:])
		delegate.ResponseWriter = w
		if err == nil {
			delegate.WriteHeader(code)
		}
		delegate.Write([]byte(SAMPLE))
	}))
)

func TestResponseWriterDelegateStatusOK(t *testing.T) {
	res, err := http.Get(ts.URL + "/200")
	assert.Nil(t, err, "No errors fetching")
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err, "No errors reading")
	assert.Equal(t, http.StatusOK, delegate.Code, "Status of 200 is recorded")
	assert.Equal(t, SAMPLE, string(body), UNCHANGED)
}

func TestResponseWriterDelegateStatusNotFound(t *testing.T) {
	res, err := http.Get(ts.URL + "/404")
	assert.Nil(t, err, "No errors fetching")
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err, "No errors reading")
	assert.Equal(t, http.StatusNotFound, delegate.Code, "Status of 404 is recorded")
	assert.Equal(t, SAMPLE, string(body), UNCHANGED)
}

func TestResponseWriterDelegateStatusUnset(t *testing.T) {
	res, err := http.Get(ts.URL + "/")
	assert.Nil(t, err, "No errors fetching")
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err, "No errors reading")
	assert.Equal(t, http.StatusOK, delegate.Code, "Status not set is 200")
	assert.Equal(t, SAMPLE, string(body), UNCHANGED)
}
