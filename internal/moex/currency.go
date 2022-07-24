package moex

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

const (
	baseURL = "https://www.banki.ru"
	path    = "/products/currency/cash"
	Region  = "moskva" // Default region
)

var lock sync.RWMutex

// Cash currency exchange rate
type Currency struct {
	name    string
	pattern string
	rate    float64
	err     error
}

func NewCurrency(name, pattern string) Currency {
	log.Println("Fetching the cash currency rate")

	lock.Lock()
	defer lock.Unlock()

	c := Currency{}
	c.name = name
	c.pattern = pattern
	c.rate, c.err = c.parseMOEX()

	return c
}

// Get formated MOEX rate
func (c *Currency) Format() string {
	lock.RLock()
	defer lock.RUnlock()

	r := fmt.Sprintf(c.pattern, c.rate)
	if c.err != nil {
		r = fmt.Sprint(c.err)
	}
	if c.rate <= 0.0 {
		r = fmt.Sprintf("%s error: wrong value of rate=%.2f", c.name, c.rate)
	}
	return r
}

// Parse MOEX
func (c *Currency) parseMOEX() (float64, error) {
	var r float64
	var err error

	s := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	s.OnHTML("div.table-flex.table-flex--no-borders.rates-summary.rates-summary--cc", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("div.table-flex__row.text-align-left.text-uppercase.font-bold.font-size-xx-large", func(i int, row *colly.HTMLElement) bool {
			c := row.ChildText("div.table-flex__cell.table-flex__cell--without-padding.padding-right-default.font-normal")
			if c != "USD" {
				return true
			}
			v := row.ChildText("div.flexbox.flexbox--vert.flexbox--gap_xxsmall.rates-summary__rate-cell")
			if len(v) == 0 {

				return true
			}
			v = strings.Replace(v, ",", ".", -1)
			r, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return true
			}
			return false
		})
	})

	s.OnError(func(r *colly.Response, err error) {
		log.Error(err)
	})

	s.Visit(baseURL + path + "/" + Region)

	return r, err
}
