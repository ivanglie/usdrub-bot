package cashex

import (
	"fmt"
	"net/http"
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
	buyBranches  map[int][]string
	sellBranches map[int][]string
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

		var url string
		if c.region != "" {
			url = baseURL + path + "/" + c.region
		} else {
			url = baseURL + path + "/" + Region
		}

		c.Lock()
		defer c.Unlock()
		c.parseBranches(url)
		c.buyMin, c.sellMin, c.buyMax, c.sellMax, c.buyAvg, c.sellAvg = findMma(c.branches)
		c.buyBranches, c.sellBranches = buyBranches(c.branches), sellBranches(c.branches)
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
func (c *Currency) BuyBranches() map[int][]string {
	c.RLock()
	defer c.RUnlock()
	return c.buyBranches
}

// SellBranches represented as string.
func (c *Currency) SellBranches() map[int][]string {
	c.RLock()
	defer c.RUnlock()
	return c.sellBranches
}

// Parse branches.
func (c *Currency) parseBranches(url string) {
	c.branches = []branch{}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	collector := colly.NewCollector()
	collector.WithTransport(t)
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

			updated, err := time.ParseInLocation("02.01.2006 15:04", row.ChildText("span.text-nowrap"), time.FixedZone("Europe/Moscow", 3))
			if err != nil {
				log.Println(err)
				return
			}

			raw := newBranch(bank, address, subway, currency, buy, sell, updated)
			if raw != (branch{}) && buy != 0 && sell != 0 && time.Now().Unix() <= updated.Local().Unix() {
				c.branches = append(c.branches, raw)
			}
		})
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Printf("UserAgent: %s", r.Headers.Get("User-Agent"))
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println(err)
	})

	err := collector.Visit(url)
	if err != nil {
		log.Println(err)
	}
}

func buyBranches(b []branch) map[int][]string {
	sort.Sort(sort.Reverse(ByBuySorter(b)))
	d := []string{}
	for i, v := range b {
		d = append(d, fmt.Sprintf("%d) %.2f RUB (_%v_): %s, %s, %s", i+1, v.Buy, v.Updated.Format("02.01.2006 15:04"), v.Bank, v.Address, v.Subway))
	}
	return func(b []string) map[int][]string {
		m := make(map[int][]string)
		j := 0
		for i := range b {
			if i%10 == 0 {
				j = i + 10

				s := []string{}
				if j < len(b) {
					s = b[i:j]
				} else {
					s = b[i:]
				}

				m[(j-10)/10] = s
			}
		}
		return m
	}(d)
}

func sellBranches(b []branch) map[int][]string {
	sort.Sort(BySellSorter(b))
	d := []string{}
	for i, v := range b {
		d = append(d, fmt.Sprintf("%d) %.2f RUB (_%v_): %s, %s, %s", i+1, v.Sell, v.Updated.Format("02.01.2006 15:04"), v.Bank, v.Address, v.Subway))
	}
	return func(b []string) map[int][]string {
		m := make(map[int][]string)
		j := 0
		for i := range b {
			if i%10 == 0 {
				j = i + 10

				s := []string{}
				if j < len(b) {
					s = b[i:j]
				} else {
					s = b[i:]
				}

				m[(j-10)/10] = s
			}
		}
		return m
	}(d)
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
