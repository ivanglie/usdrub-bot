package coingate

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getRate(t *testing.T) {
	Debug = true

	fetchFunc := func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("2"))),
		}, nil
	}

	r, err := getRate("C1", "C2", fetchFunc)
	assert.Nil(t, err)
	assert.Equal(t, float64(2), r)

	// Error from fetch
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return nil, fmt.Errorf("error")
	}
	r, err = getRate("C1", "C2", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)

	// Error from server
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 503,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"message":"Service Unavailable","reason":"ServiceUnavailable"}`))),
		}, nil
	}

	r, err = getRate("C1", "C2", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)

	// Unmarshal error
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 503,
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
		}, nil
	}

	r, err = getRate("C1", "C2", fetchFunc)
	t.Log(err)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)

	// Empty response
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
		}, nil
	}

	r, err = getRate("C1", "C2", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)

	// Not a float response
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("not a float"))),
		}, nil
	}

	r, err = getRate("C1", "C2", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)
}
