package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cbr "github.com/ivanglie/go-cbr-client"
	forex "github.com/ivanglie/go-coingate-client"
	moex "github.com/ivanglie/go-moex-client"
	"github.com/ivanglie/usdrub-bot/internal/cashex"
	"github.com/ivanglie/usdrub-bot/internal/ex"
	"github.com/ivanglie/usdrub-bot/internal/scheduler"
	"github.com/ivanglie/usdrub-bot/internal/storage"
	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

const (
	helpCmd    = "Just use /forex, /moex, /cbrf, /cash and /home command."
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

	cashKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Buy cash", "Buy"),
			tgbotapi.NewInlineKeyboardButtonData("Sell cash", "Sell"),
			tgbotapi.NewInlineKeyboardButtonData("Help", "Help"),
		),
	)

	fx   *ex.Currency
	mx   *ex.Currency
	cbrf *ex.Currency
	cash *cashex.Currency
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
	scheduler.Debug = opts.Dbg
	forex.Debug = opts.Dbg
	moex.Debug = opts.Dbg
	cbr.Debug = opts.Dbg
	cashex.Debug = opts.Dbg

	tgbotapi.SetLogger(log)
	forex.SetLogger(log)
	cbr.SetLogger(log)

	mx = ex.New(func() (float64, error) { return moex.NewClient().GetRate(moex.USDRUB) })
	fx = ex.New(func() (float64, error) { return forex.NewClient().GetRate("USD", "RUB") })
	cbrf = ex.New(func() (float64, error) { return cbr.NewClient().GetRate("USD", time.Now()) })
	cash = cashex.New(cashex.Region)

	updateRates()
	scheduler.StartCmdOnSchedule(updateRates)

	run()
}

func run() {
	bot, err := tgbotapi.NewBotAPI(opts.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = opts.Dbg
	log.Debugf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		// Check if we've gotten a message update.
		if update.Message != nil {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			if !update.Message.IsCommand() {
				continue
			}

			switch update.Message.Command() {
			case "start", "home":
				err = storage.Persist(update.Message.From)
				if err != nil {
					log.Error(err)
				}
				msg.Text = fmt.Sprintf("*%s*\n%s %s\n%s %s\n%s %s\n*%s*\n%s\n%s",
					exPrefix, fx, fxSuffix, mx, mxSuffix, cbrf, cbrfSuffix,
					cashPrefix, cash, cashSuffix)
				msg.ReplyMarkup = cashKeyboard
			case "forex":
				msg.Text = fmt.Sprintln(exPrefix, fx, fxSuffix)
			case "moex":
				msg.Text = fmt.Sprintln(exPrefix, mx, mxSuffix)
			case "cbrf":
				msg.Text = fmt.Sprintln(exPrefix, cbrf, cbrfSuffix)
			case "cash":
				msg.Text = fmt.Sprintf("%s\n%s\n%s", cashPrefix, cash, cashSuffix)
				msg.ReplyMarkup = cashKeyboard
			case "help":
				err = storage.Persist(update.Message.From)
				if err != nil {
					log.Error(err)
				}
				msg.Text = helpCmd
			default:
				msg.Text = unknownCmd
			}

			// Send the message.
			msg.ParseMode = "markdown"
			if _, err = bot.Send(msg); err != nil {
				log.Error(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err = bot.Request(callback); err != nil {
				log.Error(err)
			}

			// And finally, send a message containing the data received.
			var msg tgbotapi.MessageConfig
			switch update.CallbackQuery.Data {
			case "Buy":
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "*Buy cash*\n"+cash.BuyBranches())
			case "Sell":
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "*Sell cash*\n"+cash.SellBranches())
			case "Help":
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, helpCmd)
			default:
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			}
			msg.ParseMode = "markdown"
			if _, err := bot.Send(msg); err != nil {
				log.Error(err)
			}
		}
	}
}

func updateRates() {
	t := time.Now()
	var wg sync.WaitGroup
	fx.Update(&wg)
	mx.Update(&wg)
	cbrf.Update(&wg)
	cash.Update(&wg)
	wg.Wait()
	log.Debugln("Elapsed time:", time.Since(t))

	if _, err := fx.Rate(); err != nil {
		log.Errorf("Forex error: %v\n", err)
	}
	if _, err := mx.Rate(); err != nil {
		log.Errorf("Moscow Exchange error: %v", err)
	}
	if _, err := cbrf.Rate(); err != nil {
		log.Errorf("Russian Central Bank error: %v", err)
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
