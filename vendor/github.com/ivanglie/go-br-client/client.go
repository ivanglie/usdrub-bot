package br

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

const (
	// Example: https://www.banki.ru/products/currency/cash/usd/moskva/
	baseURL = "https://www.banki.ru/products/currency/cash/%s/%s/"

	// Currency
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	CHF Currency = "CHF"
	JPY Currency = "JPY"
	CNY Currency = "CNY"
	CZK Currency = "CZK"
	PLN Currency = "PLN"

	// City
	Barnaul         City = "barnaul"
	Voronezh        City = "voronezh"
	Volgograd       City = "volgograd"
	Vladivostok     City = "vladivostok"
	Ekaterinburg    City = "ekaterinburg"
	Irkutsk         City = "irkutsk"
	Izhevsk         City = "izhevsk"
	Kazan           City = "kazan~"
	Krasnodar       City = "krasnodar"
	Krasnoyarsk     City = "krasnoyarsk"
	Kaliningrad     City = "kaliningrad"
	Kirov           City = "kirov"
	Kemerovo        City = "kemerovo"
	Moscow          City = "moskva"
	Novosibirsk     City = "novosibirsk"
	NizhnyNovgorod  City = "nizhniy_novgorod"
	Omsk            City = "omsk"
	Orenburg        City = "orenburg"
	Perm            City = "perm~"
	RostovOnDon     City = "rostov-na-donu"
	SaintPetersburg City = "sankt-peterburg"
	Samara          City = "samara"
	Saratov         City = "saratov"
	Sochi           City = "krasnodarskiy_kray/sochi"
	Tyumen          City = "tyumen~"
	Tolyatti        City = "samarskaya_oblast~/tol~yatti"
	Tomsk           City = "tomsk"
	Ufa             City = "ufa"
	Khabarovsk      City = "habarovsk"
	Chelyabinsk     City = "chelyabinsk"
)

type Client struct {
	collector *colly.Collector
	url       string
}

func NewClient() *Client {
	c := colly.NewCollector()

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c.WithTransport(t)
	c.AllowURLRevisit = true

	extensions.RandomUserAgent(c)

	return &Client{c, ""}
}

var Debug bool

// Rates by currency (USD, if empty) and city (Moscow, if empty).
func (c *Client) Rates(crnc Currency, ct City) (*Rates, error) {
	if len(crnc) == 0 {
		crnc = USD
	}

	if len(ct) == 0 {
		ct = Moscow
	}

	if len(c.url) == 0 { // for tests
		c.url = fmt.Sprintf(baseURL, strings.ToLower(string(crnc)), ct)
	}

	if Debug {
		log.Printf("Fetching the currency rate from %s", c.url)
	}

	r := &Rates{Currency: crnc, City: ct}
	b, err := c.parseBranches()
	if err != nil {
		r = nil
	} else {
		r.Branches = b
	}

	return r, err
}

// Parse banks and their branches info.
func (c *Client) parseBranches() ([]Branch, error) {
	var b []Branch
	var err error

	c.collector.OnHTML("div.table-flex.trades-table.table-product", func(e *colly.HTMLElement) {
		e.ForEach("div.table-flex__row.item.calculator-hover-icon__container", func(i int, row *colly.HTMLElement) {
			var location *time.Location
			location, err = time.LoadLocation("Europe/Moscow")
			if err != nil {
				log.Println(err)
				return
			}

			var updated time.Time
			updated, err = time.ParseInLocation("02.01.2006 15:04", row.ChildText("span.text-nowrap"), location)
			if err != nil {
				log.Println(err)
				return
			}

			if time.Now().In(location).Sub(updated) > 24*time.Hour {
				err = fmt.Errorf("exchange rate is out of date for 24 hours: %v", updated)
				log.Println(err)
				return
			}

			bank := row.ChildText("a.font-bold")

			a := strings.Split(row.ChildAttr("a.font-bold", "data-currency-rates-tab-item"), "_")
			address := a[len(a)-1]

			subway := row.ChildText("div.font-size-small")
			currency := row.ChildAttr("div.table-flex__rate.font-size-large", "data-currencies-code")

			var buy float64
			buy, err = strconv.ParseFloat(row.ChildAttr("div.table-flex__rate.font-size-large", "data-currencies-rate-buy"), 64)
			if err != nil {
				log.Println(err)
				return
			}

			var sell float64
			sell, err = strconv.ParseFloat(row.ChildAttr("div.table-flex__rate.font-size-large.text-nowrap", "data-currencies-rate-sell"), 64)
			if err != nil {
				log.Println(err)
				return
			}

			raw := newBranch(bank, address, subway, currency, buy, sell, updated)
			if raw != (Branch{}) && (buy != 0 || sell != 0) {
				b = append(b, raw)
			}
		})
	})

	c.collector.OnRequest(func(r *colly.Request) {
		log.Printf("UserAgent: %s", r.Headers.Get("User-Agent"))
	})

	c.collector.OnError(func(r *colly.Response, e error) {
		err = e
		log.Println(err)
	})

	err = c.collector.Visit(c.url)

	return b, err
}
