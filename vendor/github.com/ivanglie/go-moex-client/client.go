package moex

import (
	"net/http"
)

// FetchFunction is a function that mimics http.Get() method
type FetchFunction func(url string) (resp *http.Response, err error)

// Client is a currency rates service client... what else?
type Client interface {
	GetRate(code string) (float64, error)
	SetFetchFunction(FetchFunction)
}

type client struct {
	fetch FetchFunction
}

func (s client) GetRate(code string) (float64, error) {
	rate, err := getRate(code, s.fetch)
	if err != nil {
		return 0, err
	}
	return rate, nil
}

func (s client) SetFetchFunction(f FetchFunction) {
	s.fetch = f
}

// NewClient creates a new rates service instance
func NewClient() Client {
	return client{http.Get}
}
