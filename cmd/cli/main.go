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
		Forex       bool `long:"forex" short:"f" description:"Exchange rate by Forex"`
		Moex        bool `long:"moex" short:"m" description:"Exchange rate by Moscow Exchange"`
		Cbrf        bool `long:"cbrf" short:"c" description:"Exchange rate by Russian Central Bank"`
		Cash        bool `long:"cash" description:"Cost of cash in Exchange Branches in Russia, Moscow"`
		CashDetails bool `long:"cashd" description:"Details of cash cost in Exchange Branches in Russia, Moscow"`
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

	mx = moex.New("1 US Dollar equals %.2f RUB by Moscow Exchange")
	forex = ex.New("1 US Dollar equals %.2f RUB by Forex", func() (float64, error) { return fx.NewClient().GetRate("USD", "RUB") })
	cbrf = ex.New("1 US Dollar equals %.2f RUB by Russian Central Bank", func() (float64, error) { return cbr.NewClient().GetRate("USD", time.Now()) })
	cash = cashex.New("1 US Dollar costs from %.2f RUB to %.2f RUB (%.2f on average) in Moscow, Russia by Banki.ru", cashex.Region)

	if opts.Forex {
		forex.Update()
		fmt.Println(forex.Format())
	}
	if opts.Moex {
		mx.Update()
		fmt.Println(mx.Format())
	}
	if opts.Cbrf {
		cbrf.Update()
		fmt.Println(cbrf.Format())
	}
	if opts.Cash {
		cash.Update()
		fmt.Println(cash.Format())
	}
	if opts.CashDetails {
		cash.Details()
		fmt.Printf("List of cash cost in Exchange Branches in Russia, Moscow:\n%s", cash.Details())
	}
}
