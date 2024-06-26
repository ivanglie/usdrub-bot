package bestchange

import (
	"net/http"
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

const (
	// Example: https://www.bestchange.com/cash-ruble-to-tether-trc20-in-msk.html.
	baseURL = "https://www.bestchange.com/cash-ruble-to-tether-trc20-in-msk.html"
)

var (
	// Debug mode. Default: false.
	Debug bool
)

// Client.
type Client struct {
	buildURL  func() string
	collector *colly.Collector
}

// NewClient creates a new client.
func NewClient() *Client {
	c := &Client{}
	c.collector = colly.NewCollector(colly.AllowURLRevisit())

	c.buildURL = func() string {
		return baseURL
	}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c.collector.WithTransport(t)
	extensions.RandomUserAgent(c.collector)

	return c
}

// Rate in Moscow.
func (c *Client) Rate() (float64, error) {
	if Debug {
		log.Printf("[DEBUG] Fetching the USDT (TRC20) rate from %s", c.buildURL())
	}

	r := &Rate{}
	b, err := c.parseRate()
	if err != nil {
		r = nil
	} else {
		r.Value = b
	}

	return r.Value, err
}

// parseRate parses rate.
func (c *Client) parseRate() (float64, error) {
	var v float64
	var err error

	c.collector.OnRequest(func(r *colly.Request) {
		if Debug {
			log.Printf("UserAgent: %s", r.Headers.Get("User-Agent"))
		}
	})

	c.collector.OnError(func(r *colly.Response, err error) {
		log.Println(err)
	})

	c.collector.OnHTML("span[title='Average rate']", func(e *colly.HTMLElement) {
		s := e.ChildText("span.bt")

		if Debug {
			log.Println("[DEBUG] Average exchange rate:", s)
		}

		if len(s) == 0 {
			return
		}

		v, err = strconv.ParseFloat(s, 64)
	})

	err = c.collector.Visit(c.buildURL())
	if err != nil {
		log.Printf("Error visiting page %v", err)
	}

	return v, err
}
