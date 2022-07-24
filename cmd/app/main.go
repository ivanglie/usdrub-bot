package main

import (
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cbr "github.com/ivanglie/go-cbr-client"
	fx "github.com/ivanglie/go-coingate-client"
	"github.com/ivanglie/usdrub-bot/internal/cashex"
	"github.com/ivanglie/usdrub-bot/internal/ex"
	"github.com/ivanglie/usdrub-bot/internal/moex"
	"github.com/ivanglie/usdrub-bot/internal/scheduler"
	"github.com/ivanglie/usdrub-bot/internal/storage"
	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

const (
	homeCmd = `Hi there!
I will help know current US Dollar to Russian Ruble exchange rate.`
	helpCmd    = "Just use /forex, /moex, /cash and /home command."
	unknownCmd = "Unknown command"

	fxTitle   = "Forex"
	cbrTitle  = "CBRF"
	moexTitle = "MOEX"
	cashTitle = "Cash"
)

var (
	homeKeyboard,
	detailsKeyboard tgbotapi.InlineKeyboardMarkup

	log *logrus.Logger

	forex,
	cbrf ex.Currency
	mx   moex.Currency
	cash cashex.Currency

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
	scheduler.Debug = opts.Dbg
	fx.Debug = opts.Dbg
	moex.Debug = opts.Dbg
	cbr.Debug = opts.Dbg
	cashex.Debug = opts.Dbg

	tgbotapi.SetLogger(log)
	fx.SetLogger(log)
	cbr.SetLogger(log)

	homeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fxTitle, fxTitle),
			tgbotapi.NewInlineKeyboardButtonData(moexTitle, moexTitle),
			tgbotapi.NewInlineKeyboardButtonData(cbrTitle, cbrTitle),
			tgbotapi.NewInlineKeyboardButtonData(cashTitle, cashTitle),
			tgbotapi.NewInlineKeyboardButtonData("Help", "Help"),
		),
	)

	detailsKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("See details", "Details"),
		),
	)

	updateRates()
	scheduler.StartCmdOnSchedule(updateRates, opts.CronSpec)

	run()
}

func updateRates() {
	var fxErr, cbrfErr error

	forex, fxErr = ex.NewCurrency(fxTitle, "1 US Dollar equals %.2f RUB by Forex", func() (float64, error) { return fx.NewClient().GetRate("USD", "RUB") })
	if fxErr != nil {
		log.Error(fxErr)
	}
	log.Debug(forex.Format())

	mx = moex.NewCurrency(moexTitle, "1 US Dollar equals %.2f RUB by Moscow Exchange")
	log.Debug(mx.Format())

	cbrf, cbrfErr = ex.NewCurrency(cbrTitle, "1 US Dollar equals %.2f RUB by Russian Central Bank", func() (float64, error) { return cbr.NewClient().GetRate("USD", time.Now()) })
	if cbrfErr != nil {
		log.Error(cbrfErr)
	}
	log.Debug(cbrf.Format())

	cash = cashex.NewCurrency(cashTitle, "1 US Dollar costs from %.2f RUB to %.2f RUB (%.2f on average) in Moscow, Russia by Banki.ru", cashex.Region)
	log.Debug(cash.Format())
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
				msg.Text = homeCmd
				msg.ReplyMarkup = homeKeyboard
			case "forex":
				msg.Text = forex.Format()
			case "moex":
				msg.Text = mx.Format()
			case "cbrf":
				msg.Text = cbrf.Format()
			case "cash":
				msg.Text = cash.Format()
				msg.ReplyMarkup = detailsKeyboard
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
			case fxTitle:
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, forex.Format())
			case moexTitle:
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, mx.Format())
			case cbrTitle:
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, cbrf.Format())
			case cashTitle:
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, cash.Format())
				msg.ReplyMarkup = detailsKeyboard
			case "Details":
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, cash.Details())
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
