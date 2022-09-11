package moex

import (
	"net/http"
)

// fetchFunction is a function that mimics http.Get() method
type fetchFunction func(url string) (resp *http.Response, err error)

// Client is a currency rates service client... what else?
type Client interface {
	GetRate(code string) (float64, error)
	SetFetchFunction(fetchFunction)
}

type client struct {
	fetch fetchFunction
}

func (s client) GetRate(code string) (float64, error) {
	rate, err := getRate(code, s.fetch)
	if err != nil {
		return 0, err
	}
	return rate, nil
}

func (s client) SetFetchFunction(f fetchFunction) {
	s.fetch = f
}

// NewClient creates a new rates service instance
func NewClient() Client {
	return client{http.Get}
}
