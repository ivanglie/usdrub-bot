package moex

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_client_GetRate(t *testing.T) {
	client := &client{}
	client.SetFetchFunction(func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(bytes.NewReader(
				[]byte(`[{"charsetinfo": {}}, {"charsetinfo": {}, "securities": [], "marketdata": [{"LAST": 2}]}]`))),
		}, nil
	})

	r, err := client.GetRate("C1")
	assert.Nil(t, err)
	assert.Equal(t, float64(2), r)

	// Error from fetch
	client.SetFetchFunction(func(url string) (resp *http.Response, err error) {
		return nil, fmt.Errorf("error")
	})

	r, err = client.GetRate("C1")
	assert.Error(t, err)
	assert.Equal(t, float64(0), r)
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)
}
