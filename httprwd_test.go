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
		code, _ := strconv.Atoi(r.URL.Path[1:])
		delegate.ResponseWriter = w
		delegate.WriteHeader(code)
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
	res, err := http.Get(ts.URL + "/0")
	// If golang is older than https://github.com/golang/go/commit/46069bed06551337bc8c9b293040d30e41917289
	if err == nil {
		body, err = ioutil.ReadAll(res.Body)
		assert.Nil(t, err, "No errors reading")
		assert.Equal(t, http.StatusOK, delegate.Code, "Status of 0 is 200")
		assert.Equal(t, SAMPLE, string(body), UNCHANGED)
	}
}
