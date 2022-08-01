package main

import (
	"fmt"
	"log"
	"os"
	"time"

	cbr "github.com/ivanglie/go-cbr-client"
	fx "github.com/ivanglie/go-coingate-client"
	"github.com/ivanglie/usdrub-bot/internal/cashex"
	"github.com/ivanglie/usdrub-bot/internal/ex"
	"github.com/ivanglie/usdrub-bot/internal/moex"
	"github.com/jessevdk/go-flags"
)

var (
	forex *ex.Currency
	cbrf  *ex.Currency
	mx    *moex.Currency
	cash  *cashex.Currency

	opts struct {
		Forex        bool `long:"forex" short:"f" description:"Exchange rate by Forex"`
		Moex         bool `long:"moex" short:"m" description:"Exchange rate by Moscow Exchange"`
		Cbrf         bool `long:"cbrf" short:"c" description:"Exchange rate by Russian Central Bank"`
		Cash         bool `long:"cash" description:"Cost of cash in Exchange Branches in Russia, Moscow"`
		CashBranches bool `long:"cashb" description:"Cash Exchange branches os Moscow, Russia"`
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
		forex.Update()
		fmt.Println(forex.Rate())
	}
	if opts.Moex {
		mx.Update()
		fmt.Println(mx.Rate())
	}
	if opts.Cbrf {
		cbrf.Update()
		fmt.Println(cbrf.Rate())
	}
	if opts.Cash {
		cash.Update()
		fmt.Println(cash.Rate())
	}
	if opts.CashBranches {
		cash.Update()
		fmt.Printf("Cash exchange rates in branches in Moscow, Russia\nBuy cash\n%sSell cash\n%s",
			cash.BuyBranches(), cash.SellBranches())
	}
}
