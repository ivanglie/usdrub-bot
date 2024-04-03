package cbr

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockHttpClient is a mock http client with status code 500.
type mockHttpClient struct{}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 500}, nil
}

// mockHttpClientErr is a mock http client with error.
type mockHttpClientErr struct{}

func (m *mockHttpClientErr) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("error")
}

// mockReadCloser is a mock read closer with error.
type mockReadCloser struct{}

func (r *mockReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (r *mockReadCloser) Close() (err error) {
	return errors.New("close error")
}

// mockHttpClientResponseBodyErr is a mock http client with response body error.
// It returns a response with status code 200 and a read closer with error.
type mockHttpClientResponseBodyErr struct{}

func (m *mockHttpClientResponseBodyErr) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: &mockReadCloser{}}, nil
}

func TestClient_GetRate(t *testing.T) {
	Debug = true

	client := NewClient()
	rate, err := client.GetRate("USD", time.Now())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, rate, float64(1))

	Debug = false

	// unknown currency: _
	rate, err = client.GetRate("_", time.Now())
	assert.Error(t, err)
	assert.Equal(t, "unknown currency: _", err.Error())
	assert.Equal(t, float64(0), rate)

	// status code: 500
	client.httpClient = &mockHttpClient{}
	rate, err = client.GetRate("CNY", time.Now())
	assert.Error(t, err)
	assert.Equal(t, "status code: 500", err.Error())
	assert.Equal(t, float64(0), rate)

	// error
	client.httpClient = &mockHttpClientErr{}
	rate, err = client.GetRate("CNY", time.Now())
	assert.Error(t, err)
	assert.Equal(t, "error", err.Error())
	assert.Equal(t, float64(0), rate)

	// response body error
	client.httpClient = &mockHttpClientResponseBodyErr{}
	rate, err = client.GetRate("CNY", time.Now())
	assert.Error(t, err)
	t.Log(err)
	assert.Equal(t, float64(0), rate)
}

func Test_currencyRateValue_Error(t *testing.T) {
	c := Currency{}
	c.Value = "0'1"
	rate, err := currencyRateValue(c)
	assert.Error(t, err)
	assert.Equal(t, float64(0), rate)
}
