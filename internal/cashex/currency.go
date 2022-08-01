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
	region       string
	branches     []branch
	buyBranches  string
	sellBranches string
	buyMin       float64
	buyMax       float64
	buyAvg       float64
	sellMin      float64
	sellMax      float64
	sellAvg      float64
}

func New(region string) *Currency {
	return &Currency{
		region: region,
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
	c.buyMin, c.sellMin, c.buyMax, c.sellMax, c.buyAvg, c.sellAvg = findMma(b)
	c.buyBranches = buyBranches(b)
	c.sellBranches = sellBranches(b)
}

// Get formated cash exchange rate: buyMin, buyMax, buyAvg, sellMin, sellMax, sellAvg
func (c *Currency) Rate() (float64, float64, float64, float64, float64, float64) {
	c.RLock()
	defer c.RUnlock()
	return c.buyMin, c.buyMax, c.buyAvg, c.sellMin, c.sellMax, c.sellAvg
}

// Get buy branches
func (c *Currency) BuyBranches() string {
	c.RLock()
	defer c.RUnlock()
	return c.buyBranches
}

// Get sell branches
func (c *Currency) SellBranches() string {
	c.RLock()
	defer c.RUnlock()
	return c.sellBranches
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

func buyBranches(b []branch) string {
	sort.Sort(sort.Reverse(ByBuySorter(ByBuySorter(b))))
	d := ""
	for i, v := range b {
		d = d + fmt.Sprintf("%d. %.2f RUB: %s, %s, %s\n", i+1, v.Buy, v.Bank, v.Address, v.Subway)
	}
	return d
}

func sellBranches(b []branch) string {
	sort.Sort(BySellSorter(b))
	d := ""
	for i, v := range b {
		d = d + fmt.Sprintf("%d. %.2f RUB: %s, %s, %s\n", i+1, v.Sell, v.Bank, v.Address, v.Subway)
	}
	return d
}

// Find min, max and avg
func findMma(b []branch) (float64, float64, float64, float64, float64, float64) {
	if len(b) == 0 {
		return 0, 0, 0, 0, 0, 0
	}
	bmin := b[0].Buy
	smin := b[0].Sell
	bmax := b[0].Buy
	smax := b[0].Sell
	btotal := 0.0
	stotal := 0.0
	for _, v := range b {
		if v.Buy < bmin {
			bmin = v.Buy
		}
		if v.Sell < smin {
			smin = v.Sell
		}
		if v.Buy > bmax {
			bmax = v.Buy
		}
		if v.Sell > smax {
			smax = v.Sell
		}
		btotal += v.Buy
		stotal += v.Sell
	}
	return bmin, smin, bmax, smax, btotal / float64(len(b)), stotal / float64(len(b))
}
