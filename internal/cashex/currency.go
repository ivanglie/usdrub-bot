package cashex

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

const (
	baseURL = "https://www.banki.ru"
	path    = "/products/currency/cash"
	Region  = "moskva" // Default region
)

var Debug bool

// Currency exchange rate of cash.
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
	return &Currency{region: region}
}

// Update currency exchange cash rate.
func (c *Currency) Update(wg *sync.WaitGroup) {
	update := func() {
		if Debug {
			log.Println("Fetching the currency rate")
		}

		c.Lock()
		defer c.Unlock()
		b := c.parseBranches(c.region)
		c.branches = b
		c.buyMin, c.sellMin, c.buyMax, c.sellMax, c.buyAvg, c.sellAvg = findMma(b)
		c.buyBranches, c.sellBranches = buyBranches(b), sellBranches(b)
	}

	if wg == nil {
		update()
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		update()
	}()
}

// Rate of currency exchange cash returns of buyMin, buyMax, buyAvg, sellMin, sellMax, sellAvg.
func (c *Currency) Rate() (float64, float64, float64, float64, float64, float64) {
	c.RLock()
	defer c.RUnlock()
	return c.buyMin, c.buyMax, c.buyAvg, c.sellMin, c.sellMax, c.sellAvg
}

// String representation of currency exchange cash rate.
func (c *Currency) String() string {
	c.RLock()
	defer c.RUnlock()
	return fmt.Sprintf("Buy:\t%.2f .. %.2f RUB (avg %.2f)\nSell:\t%.2f .. %.2f RUB (avg %.2f)",
		c.buyMax, c.buyMin, c.buyAvg, c.sellMin, c.sellMax, c.sellAvg)
}

// BuyBranches represented as string.
func (c *Currency) BuyBranches() string {
	c.RLock()
	defer c.RUnlock()
	return c.buyBranches
}

// SellBranches represented as string.
func (c *Currency) SellBranches() string {
	c.RLock()
	defer c.RUnlock()
	return c.sellBranches
}

// Parse branches.
func (c *Currency) parseBranches(region string) (branches []branch) {
	if len(region) == 0 {
		region = Region
	}

	collector := colly.NewCollector()
	collector.AllowURLRevisit = true
	extensions.RandomUserAgent(collector)

	collector.OnHTML("div.table-flex.trades-table.table-product", func(e *colly.HTMLElement) {
		e.ForEach("div.table-flex__row.item.calculator-hover-icon__container", func(i int, row *colly.HTMLElement) {
			bank := row.ChildText("a.font-bold")

			a := strings.Split(row.ChildAttr("a.font-bold", "data-currency-rates-tab-item"), "_")
			address := a[len(a)-1]

			subway := row.ChildText("div.font-size-small")
			currency := row.ChildAttr("div.table-flex__rate.font-size-large", "data-currencies-code")

			buy, err := strconv.ParseFloat(row.ChildAttr("div.table-flex__rate.font-size-large", "data-currencies-rate-buy"), 64)
			if err != nil {
				log.Println(err)
				return
			}

			sell, err := strconv.ParseFloat(row.ChildAttr("div.table-flex__rate.font-size-large.text-nowrap", "data-currencies-rate-sell"), 64)
			if err != nil {
				log.Println(err)
				return
			}

			updated, _ := time.Parse("02.01.2006 15:04", row.ChildText("span.text-nowrap"))
			if updated.YearDay() < time.Now().YearDay() {
				return
			}

			branches = append(branches, newBranch(bank, address, subway, currency, buy, sell, updated))
		})
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Printf("UserAgent: %s", r.Headers.Get("User-Agent"))
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println(err)
	})

	err := collector.Visit(baseURL + path + "/" + region)
	if err != nil {
		log.Println(err)
	}

	return
}

func buyBranches(b []branch) string {
	sort.Sort(sort.Reverse(ByBuySorter(b)))
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
func findMma(b []branch) (bmin, smin, bmax, smax, bavg, savg float64) {
	if len(b) == 0 {
		return
	}

	bmin, smin, bmax, smax = b[0].Buy, b[0].Sell, b[0].Buy, b[0].Sell
	btotal, stotal := float64(0), float64(0)
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
	bavg, savg = btotal/float64(len(b)), stotal/float64(len(b))

	return
}
