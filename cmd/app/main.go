package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	br "github.com/ivanglie/go-br-client"
	cbr "github.com/ivanglie/go-cbr-client"
	forex "github.com/ivanglie/go-coingate-client"
	moex "github.com/ivanglie/go-moex-client"
	"github.com/ivanglie/usdrub-bot/internal/exrate"
	"github.com/ivanglie/usdrub-bot/internal/scheduler"
	"github.com/ivanglie/usdrub-bot/internal/storage"
	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

const (
	helpCmd    = "Just use /forex, /moex, /cbrf, /cash and /dashboard command."
	unknownCmd = "Unknown command"
	exPrefix   = "1 US Dollar equals"
	cashPrefix = "Cash exchange rates"
	fxSuffix   = "by Forex"
	mxSuffix   = "by Moscow Exchange"
	cbrfSuffix = "by Russian Central Bank"
	cashSuffix = "in branches in Moscow, Russia by Banki.ru"
)

var (
	log *logrus.Logger

	opts struct {
		Dbg      bool   `long:"dbg" env:"DEBUG" description:"Debug mode"`
		BotToken string `long:"bottoken" env:"BOT_TOKEN" description:"Telegram API Token"`
		CronSpec string `long:"cronspec" env:"CRON_SPEC" description:"Cron spec"`
	}

	version = "unknown"

	bot *tgbotapi.BotAPI

	cashKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Buy cash", "Buy"),
			tgbotapi.NewInlineKeyboardButtonData("Sell cash", "Sell"),
			tgbotapi.NewInlineKeyboardButtonData("Help", "Help"),
		),
	)

	cashBuyMoreKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("More", "BuyMore"),
		),
	)

	cashSellMoreKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("More", "SellMore"),
		),
	)

	fx, mx, cbrf *exrate.Rate
	cash         *exrate.CashRate
	cbb, csb     map[int64]int // Current Buy Branches (cbb) and Sell Branches (csb) for chat ID
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
	scheduler.Debug, forex.Debug, moex.Debug, cbr.Debug, br.Debug = opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg

	tgbotapi.SetLogger(log)
	scheduler.SetLogger(log)
	forex.SetLogger(log)
	cbr.SetLogger(log)
	moex.SetLogger(log)
	br.SetLogger(log)

	mx = exrate.NewRate(func() (float64, error) { return moex.NewClient().GetRate(moex.USDRUB) })
	fx = exrate.NewRate(func() (float64, error) { return forex.NewClient().GetRate("USD", "RUB") })
	cbrf = exrate.NewRate(func() (float64, error) { return cbr.NewClient().GetRate("USD", time.Now()) })
	cash = exrate.NewCashRate(func() (*br.Rates, error) { return br.NewClient().Rates(br.USD, br.Moscow) })

	updateRates()
	scheduler.StartCmdOnSchedule(updateRates)

	run()
}

func run() {
	var err error
	bot, err = tgbotapi.NewBotAPI(opts.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = opts.Dbg
	log.Debugf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		msg := tgbotapi.MessageConfig{}
		if update.Message != nil {
			if !update.Message.IsCommand() {
				continue
			}

			switch update.Message.Command() {
			case "start", "dashboard", "help":
				err := storage.Persist(update.Message.From)
				if err != nil {
					log.Error(err)
				}
			}

			msg = messageByCommand(update.Message.Chat.ID, update.Message.Command())
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Error(err)
			}

			msg = messageByCallbackData(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		}

		msg.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(msg); err != nil {
			log.Error(err)
		}
	}
}

func messageByCommand(chatId int64, command string) (m tgbotapi.MessageConfig) {
	m.ChatID = chatId

	switch command {
	case "start", "dashboard":
		m.Text = fmt.Sprintf("*%s*\n%s %s\n%s %s\n%s %s\n*%s*\n%s\n%s",
			exPrefix, fx, fxSuffix, mx, mxSuffix, cbrf, cbrfSuffix, cashPrefix, cash, cashSuffix)
		m.ReplyMarkup = cashKeyboard
	case "forex":
		m.Text = fmt.Sprintln(exPrefix, fx, fxSuffix)
	case "moex":
		m.Text = fmt.Sprintln(exPrefix, mx, mxSuffix)
	case "cbrf":
		m.Text = fmt.Sprintln(exPrefix, cbrf, cbrfSuffix)
	case "cash":
		m.Text = fmt.Sprintf("%s\n%s\n%s", cashPrefix, cash, cashSuffix)
		m.ReplyMarkup = cashKeyboard
	case "help":
		m.Text = helpCmd
	default:
		m.Text = unknownCmd
	}

	return
}

func messageByCallbackData(chatId int64, data string) (m tgbotapi.MessageConfig) {
	m.ChatID = chatId

	switch data {
	case "Buy":
		b := cash.BuyBranches()
		cbb = make(map[int64]int)
		cbb[chatId] = 0
		m.Text = "*Buy cash*\n" + strings.Join(b[cbb[chatId]], "\n")
		if len(b) > 1 {
			m.ReplyMarkup = cashBuyMoreKeyboard
		}
	case "BuyMore":
		b := cash.BuyBranches()
		if cbb[chatId] < len(b) {
			cbb[chatId] = cbb[chatId] + 1
		}

		if b[cbb[chatId]] != nil {
			m.Text = "*Buy cash*\n" + strings.Join(b[cbb[chatId]], "\n")
		}

		if cbb[chatId] != len(b)-1 {
			m.ReplyMarkup = cashBuyMoreKeyboard
		}
	case "Sell":
		b := cash.SellBranches()
		csb = make(map[int64]int)
		csb[chatId] = 0
		m.Text = "*Sell cash*\n" + strings.Join(b[csb[chatId]], "\n")
		if len(b) > 1 {
			m.ReplyMarkup = cashSellMoreKeyboard
		}
	case "SellMore":
		b := cash.SellBranches()
		if csb[chatId] < len(b) {
			csb[chatId] = csb[chatId] + 1
		}

		if b[csb[chatId]] != nil {
			m.Text = "*Sell cash*\n" + strings.Join(b[csb[chatId]], "\n")
		}

		if csb[chatId] != len(b)-1 {
			m.ReplyMarkup = cashSellMoreKeyboard
		}
	case "Help":
		m.Text = helpCmd
	default:
		m.Text = data
	}

	return
}

func updateRates() {
	t := time.Now()
	wg := &sync.WaitGroup{}
	fx.Update(wg)
	mx.Update(wg)
	cbrf.Update(wg)
	cash.Update(wg)
	wg.Wait()
	log.Debugln("Elapsed time:", time.Since(t))

	if _, err := fx.Rate(); err != nil {
		log.Errorf("error by Forex: %v\n", err)
	}
	if _, err := mx.Rate(); err != nil {
		log.Errorf("error by Moscow Exchange: %v", err)
	}
	if _, err := cbrf.Rate(); err != nil {
		log.Errorf("error by Russian Central Bank: %v", err)
	}
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
