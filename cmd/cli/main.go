package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	cbr "github.com/ivanglie/go-cbr-client"
	fx "github.com/ivanglie/go-coingate-client"
	"github.com/ivanglie/usdrub-bot/internal/cashex"
	"github.com/ivanglie/usdrub-bot/internal/ex"
	"github.com/ivanglie/usdrub-bot/internal/moex"
	"github.com/jessevdk/go-flags"
	"github.com/olekukonko/tablewriter"
)

var (
	forex *ex.Currency
	cbrf  *ex.Currency
	mx    *moex.Currency
	cash  *cashex.Currency

	opts struct {
		Forex    bool `long:"forex" description:"Exchange rate by Forex"`
		Moex     bool `long:"moex" description:"Exchange rate by Moscow Exchange"`
		Cbrf     bool `long:"cbrf" description:"Exchange rate by Russian Central Bank"`
		Cash     bool `long:"cash" description:"Cash exchange rates in Russia, Moscow"`
		BuyCash  bool `long:"buy" description:"Buy cash"`
		SellCash bool `long:"sell" description:"Sell cash"`
		Rates    bool `long:"rates" description:"Exchange rate by Forex, Moscow Exchange and Russian Central Bank"`
	}
)

func main() {
	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Printf("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}

	mx = moex.New()
	forex = ex.New(func() (float64, error) { return fx.NewClient().GetRate("USD", "RUB") })
	cbrf = ex.New(func() (float64, error) { return cbr.NewClient().GetRate("USD", time.Now()) })
	cash = cashex.New(cashex.Region)

	if opts.Forex {
		forex.Update(nil)
		fmt.Println("1 US Dollar equals", forex, "by Forex")
	}
	if opts.Moex {
		mx.Update(nil)
		fmt.Println("1 US Dollar equals", mx, "by Moscow Exchange")
	}
	if opts.Cbrf {
		cbrf.Update(nil)
		fmt.Println("1 US Dollar equals", cbrf, "by Russian Central Bank")
	}
	if opts.Cash {
		cash.Update(nil)
		bmn, bmx, ba, smn, smx, sa := cash.Rate()
		rates := [][]string{
			{"Min", fmt.Sprintf("%.2f", bmn), fmt.Sprintf("%.2f", smn)},
			{"Max", fmt.Sprintf("%.2f", bmx), fmt.Sprintf("%.2f", smx)},
			{"Avg", fmt.Sprintf("%.2f", ba), fmt.Sprintf("%.2f", sa)},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"", "Buy, RUB", "Sell, RUB"})

		for _, v := range rates {
			table.Append(v)
		}
		table.Render() // Send output
	}
	if opts.Rates {
		t := time.Now()
		var wg sync.WaitGroup
		forex.Update(&wg)
		mx.Update(&wg)
		cbrf.Update(&wg)
		wg.Wait()
		log.Println("Elapsed time:", time.Since(t))
		rates := [][]string{
			{forex.String(), "by Forex"},
			{mx.String(), "by Moscow Exchange"},
			{cbrf.String(), "by Russian Central Bank"},
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"1 US Dollar equals", "Source"})

		for _, v := range rates {
			table.Append(v)
		}
		table.Render() // Send output
	}
	if opts.BuyCash {
		branches("Buy")
	}
	if opts.SellCash {
		branches("Sell")
	}
}

func branches(title string) {
	cash.Update(nil)
	var b string
	switch title {
	case "Buy":
		b = cash.BuyBranches()
	case "Sell":
		b = cash.SellBranches()
	default:
		b = ""
	}
	func(branches string) {
		rates := [][]string{
			{branches},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{title})
		table.SetAutoWrapText(false)

		for _, v := range rates {
			table.Append(v)
		}
		table.Render() // Send output
	}(b)
}
