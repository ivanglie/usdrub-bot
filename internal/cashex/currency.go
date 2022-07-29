package cashex

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

const (
	baseURL = "https://www.banki.ru"
	path    = "/products/currency/cash"
	Region  = "moskva" // Default region
)

var (
	Debug bool
)

// Cash currency exchange rate
type Currency struct {
	sync.RWMutex
	pattern  string
	region   string
	details  string
	branches []branch
	min      float64
	max      float64
	avg      float64
}

func New(pattern, region string) *Currency {
	return &Currency{
		pattern: pattern,
		region:  region,
	}
}

func (c *Currency) Update() {
	if Debug {
		log.Println("Fetching the currency rate")
	}
	c.Lock()
	defer c.Unlock()
	b := c.parseBranches(c.region)
	c.branches = b
	c.details = details(b)
	c.min = min(b)
	c.max = max(b)
	c.avg = avg(b)
}

// Get formated cash exchange rate
func (c *Currency) Format() string {
	c.RLock()
	defer c.RUnlock()
	r := fmt.Sprintf(c.pattern, c.min, c.max, c.avg)
	if c.min <= 0.0 || c.max <= 0.0 || c.avg <= 0.0 {
		r = fmt.Sprintf("cashex error: wrong value of min=%.2f, max=%.2f, avg=%.2f", c.min, c.max, c.avg)
	}
	return r
}

// Get details
func (c *Currency) Details() string {
	c.RLock()
	defer c.RUnlock()
	return c.details
}

// Parse branches
func (c *Currency) parseBranches(region string) []branch {
	var branches []branch

	if len(region) == 0 {
		region = Region
	}

	s := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	s.OnHTML("div.table-flex.trades-table.table-product", func(e *colly.HTMLElement) {
		e.ForEach("div.table-flex__row.item.calculator-hover-icon__container", func(i int, row *colly.HTMLElement) {
			bank := row.ChildText("a.font-bold")

			a := strings.Split(row.ChildAttr("a.font-bold", "data-currency-rates-tab-item"), "_")
			address := a[len(a)-1]

			subway := row.ChildText("div.font-size-small")
			currency := row.ChildAttr("div.table-flex__rate.font-size-large", "data-currencies-code")

			buy, err := strconv.ParseFloat(row.ChildAttr("div.table-flex__rate.font-size-large", "data-currencies-rate-buy"), 64)
			if err != nil {
				return
			}

			sell, err := strconv.ParseFloat(row.ChildAttr("div.table-flex__rate.font-size-large.text-nowrap", "data-currencies-rate-sell"), 64)
			if err != nil {
				return
			}

			updated, _ := time.Parse("02.01.2006 15:04", row.ChildText("span.text-nowrap"))
			if updated.YearDay() < time.Now().YearDay() {
				return
			}

			branches = append(branches, newBranch(bank, address, subway, currency, buy, sell, updated))
		})
	})

	s.OnError(func(r *colly.Response, err error) {
		log.Error(err)
	})

	s.Visit(baseURL + path + "/" + region)

	return branches
}

func details(b []branch) string {
	sort.Sort(BySellSorter(b))
	details := ""
	for i, v := range b {
		details = details + fmt.Sprintf("%d. *%.2f* RUB: %s, %s, %s\n", i+1, v.Sell, v.Bank, v.Address, v.Subway)
	}
	return details
}

func min(b []branch) float64 {
	if len(b) == 0 {
		return 0
	}

	min := b[0].Sell
	for _, v := range b {
		if v.Sell < min {
			min = v.Sell
		}
	}
	return min
}

func max(b []branch) float64 {
	if len(b) == 0 {
		return 0
	}

	max := b[0].Sell
	for _, v := range b {
		if v.Sell > max {
			max = v.Sell
		}
	}
	return max
}

func avg(b []branch) float64 {
	if len(b) == 0 {
		return 0
	}

	total := 0.0
	for _, v := range b {
		total += v.Sell
	}
	return total / float64(len(b))
}
