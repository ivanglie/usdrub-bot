package moex

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
			Body: io.NopCloser(bytes.NewReader(
				[]byte(`[{"charsetinfo": {}}, {"charsetinfo": {}, "securities": [], "marketdata": [{"LAST": 2}]}]`))),
		}, nil
	}

	r, err := getRate("C1", fetchFunc)
	assert.Nil(t, err)
	assert.Equal(t, float64(2), r)

	// Error from fetch
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return nil, fmt.Errorf("error")
	}
	r, err = getRate("C1", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)

	// Unmarshal error
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
		}, nil
	}

	r, err = getRate("C1", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)

	// length of c.values less than 2
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[{}]`))),
		}, nil
	}

	r, err = getRate("C1", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)

	// val.Marketdata is zero
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[{"charsetinfo": {}}, {"charsetinfo": {}, "securities": []}]`))),
		}, nil
	}

	r, err = getRate("C1", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)

	// Length of md equals 0
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[{"charsetinfo": {}}, {"charsetinfo": {}, "securities": [], "marketdata": []}]`))),
		}, nil
	}

	r, err = getRate("C1", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)

	// Not a float response
	fetchFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`[{"charsetinfo": {}}, {"charsetinfo": {}, "securities": [], "marketdata": [{"LAST": "2"}]}]`))),
		}, nil
	}

	r, err = getRate("C1", fetchFunc)
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)
}

func TestCurrency_String(t *testing.T) {
	s := Currency{
		values: []currency{
			{
				Charsetinfo: struct {
					Name string "json:\"name\""
				}{
					Name: "UTF-8",
				},
			},
		},
	}

	assert.NotNil(t, s.String())

	// Empty values
	s.values = []currency{}
	assert.Empty(t, s.String())
}
