package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ivanglie/go-br-client"
	"github.com/ivanglie/go-cbr-client"
	"github.com/ivanglie/go-coingate-client"
	"github.com/ivanglie/go-moex-client"
	"github.com/ivanglie/usdrub-bot/internal/bot"
	"github.com/ivanglie/usdrub-bot/internal/cexrate"
	"github.com/ivanglie/usdrub-bot/internal/exrate"
	"github.com/ivanglie/usdrub-bot/internal/utils"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger

	opts struct {
		Dbg      bool   `long:"dbg" env:"DEBUG" description:"Debug mode"`
		BotToken string `long:"bottoken" env:"BOT_TOKEN" description:"Telegram API Token"`
		CronSpec string `long:"cronspec" env:"CRON_SPEC" description:"Cron spec"`
	}

	version = "unknown"
)

func main() {
	fmt.Printf("usdrub-bot %s\n", version)

	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Printf("[ERROR] usdrub-bot error: %v", err)
		}
		os.Exit(2)
	}

	setupLog(opts.Dbg)
	setLogger(log)
	coingate.Debug, moex.Debug, cbr.Debug, br.Debug, utils.Debug, bot.Debug = opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg

	updateRates := func() {
		t := time.Now()

		type RateInterface interface {
			Update()
		}

		rates := []RateInterface{exrate.Get(), cexrate.Get()}

		wg := sync.WaitGroup{}
		for _, r := range rates {
			wg.Add(1)
			go func(r RateInterface) {
				defer wg.Done()
				r.Update()
			}(r)
		}

		wg.Wait()
		log.Debugln("Elapsed time:", time.Since(t))
	}

	updateRates()

	if err := utils.StartCmdOnSchedule(updateRates); err != nil {
		log.Panic(err)
	}

	bot.CreateAndStart(opts.BotToken)
}

func setupLog(dbg bool) {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		DisableLevelTruncation: true,
		FullTimestamp:          true,
		TimestampFormat:        time.RFC3339,
	})
	if dbg {
		log.SetLevel(logrus.DebugLevel)
		return
	}
	log.SetLevel(logrus.ErrorLevel)
}

func setLogger(log *logrus.Logger) {
	coingate.SetLogger(log)
	cbr.SetLogger(log)
	moex.SetLogger(log)
	br.SetLogger(log)
	utils.SetLogger(log)
	bot.SetLogger(log)
}
