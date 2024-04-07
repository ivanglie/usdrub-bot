package coingate

import (
	"net/http"
)

// FetchFunction is a function that mimics http.Get() method.
type FetchFunction func(url string) (resp *http.Response, err error)

// Client is the interface for the rates service.
type Client interface {
	GetRate(from, to string) (float64, error)
	SetFetchFunction(FetchFunction)
}

type client struct {
	fetch FetchFunction
}

// GetRate returns the exchange rate between two currencies.
// Arguments are ISO Symbol. Example: EUR, USD, BTC, ETH, etc.
// See https://developer.coingate.com/docs/get-rate
func (s *client) GetRate(from, to string) (float64, error) {
	rate, err := getRate(from, to, s.fetch)
	if err != nil {
		return 0, err
	}
	return rate, nil
}

// SetFetchFunction allows to set a custom fetch function.
func (s *client) SetFetchFunction(f FetchFunction) {
	s.fetch = f
}

// NewClient creates a new rates service instance.
func NewClient() Client {
	return &client{http.Get}
}
