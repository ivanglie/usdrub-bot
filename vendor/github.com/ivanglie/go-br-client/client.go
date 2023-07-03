package br

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

const (
	// Example: https://www.banki.ru/products/currency/map/moskva/.
	baseURL = "https://www.banki.ru/products/currency/map/%s/"

	// Currency.
	currency = "USD"

	// City.
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

var (
	// Debug mode. Default: false.
	Debug bool
)

// Client.
type Client struct {
	city      City
	buildURL  func() string
	collector *colly.Collector
}

// NewClient creates a new client.
func NewClient() *Client {
	c := &Client{}

	c.city = Moscow
	c.buildURL = func() string {
		return fmt.Sprintf(baseURL, c.city)
	}
	c.collector = colly.NewCollector(colly.AllowURLRevisit())

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c.collector.WithTransport(t)
	extensions.RandomUserAgent(c.collector)

	return c
}

// Rates USDRUB by and city (Moscow, if empty).
func (c *Client) Rates(ct City) (*Rates, error) {
	if len(ct) > 0 {
		c.city = ct
	}

	if Debug {
		log.Printf("[DEBUG] Fetching the currency rate from %s", c.buildURL())
	}

	r := &Rates{Currency: currency, City: ct}
	b, err := c.parseBranches()
	if err != nil {
		r = nil
	} else {
		r.Branches = b
	}

	if Debug {
		log.Printf("[DEBUG] Found %d branches", len(b))
	}

	return r, err
}

// parseBranches parses branches info.
func (c *Client) parseBranches() ([]Branch, error) {
	var b []Branch
	var err error

	c.collector.OnRequest(func(r *colly.Request) {
		if Debug {
			log.Printf("UserAgent: %s", r.Headers.Get("User-Agent"))
		}
	})

	c.collector.OnError(func(r *colly.Response, err error) {
		log.Println(err)
	})

	c.collector.OnHTML(".fdpae", func(e *colly.HTMLElement) {
		e.ForEach(".cITBmP", func(i int, row *colly.HTMLElement) {
			raw, err := parseBranch(row)
			if raw != (Branch{}) && err == nil {
				b = append(b, raw)
			}
		})
	})

	err = c.collector.Visit(c.buildURL())
	if err != nil {
		log.Printf("Error visiting page %v", err)
	}

	return b, err
}

// parseBranch parses branch info from the HTML element.
func parseBranch(e *colly.HTMLElement) (Branch, error) {
	bank := sanitaze(e.ChildText(".dPnGDN"))

	updatedDateString := sanitaze(e.ChildText(".cURBaH"))
	if len(updatedDateString) == 0 {
		updatedDateString = sanitaze(e.ChildText(".kiIzbr"))
	}

	s := strings.Split(updatedDateString, " ")
	if count := len(s); count >= 3 {
		updatedDateString = strings.Join(s[count-2:], " ")
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return Branch{}, err
	}

	updatedDate, err := time.ParseInLocation("02.01.2006 15:04", updatedDateString, loc)
	if err != nil {
		return Branch{}, err
	}

	if time.Now().In(loc).Sub(updatedDate) > 24*time.Hour {
		return Branch{}, fmt.Errorf("exchange rate is out of date for 24 hours: %v", updatedDate)
	}

	ratesString := sanitaze(e.ChildText(".fvORFF"))
	if len(ratesString) == 0 {
		ratesString = sanitaze(e.ChildText(".guOAzm"))
	}

	var buyRateString, sellRateString string
	if s := strings.Split(ratesString, "₽"); len(s) >= 2 {
		buyRateString = s[0]
		sellRateString = s[1]
	}

	buyRateString = strings.Replace(buyRateString, " ", "", -1)
	buyRate, err := strconv.ParseFloat(strings.ReplaceAll(buyRateString, ",", "."), 64)
	if err != nil || buyRate <= 0 {
		return Branch{}, err
	}

	sellRateString = strings.Replace(sellRateString, " ", "", -1)
	sellRate, err := strconv.ParseFloat(strings.ReplaceAll(sellRateString, ",", "."), 64)
	if err != nil || sellRate <= 0 {
		return Branch{}, err
	}

	subway := sanitaze(e.ChildText(".eybsgm"))

	return newBranch(bank, subway, currency, buyRate, sellRate, updatedDate), nil
}

// sanitaize string.
func sanitaze(s string) string {
	if len(s) == 0 {
		return s
	}

	s = strings.Replace(s, "\n", "", -1)
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")

	return s
}
