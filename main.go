package main

import (
	"fmt"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cbr "github.com/ivanglie/go-cbr-client"
	fx "github.com/ivanglie/go-coingate-client"
	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

const (
	homeCmd = `Hi there!
I will help know current US Dollar to Russian Ruble exchange rate.`
	helpCmd    = "Just use /forex, /cbrf and /home command."
	unknownCmd = "Unknown command"

	fxTitle  = "Forex"
	cbrTitle = "CBRF"
)

var (
	homeKeyboard tgbotapi.InlineKeyboardMarkup

	log *logrus.Logger

	forex,
	cbrf Source

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
			log.Printf("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}

	setupLog(opts.Dbg)
	fx.Debug = opts.Dbg
	cbr.Debug = opts.Dbg

	tgbotapi.SetLogger(log)
	fx.SetLogger(log)
	cbr.SetLogger(log)

	homeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fxTitle, fxTitle),
			tgbotapi.NewInlineKeyboardButtonData("CBRF", cbrTitle),
			tgbotapi.NewInlineKeyboardButtonData("Help", "Help"),
		),
	)

	forex = Source{
		name:     fxTitle,
		pattern:  "1 US Dollar equals %.2f RUB by Forex",
		rateFunc: func() (float64, error) { return fx.NewClient().GetRate("USD", "RUB") },
	}

	cbrf = Source{
		name:     cbrTitle,
		pattern:  "1 US Dollar equals %.2f RUB by Russian Central Bank",
		rateFunc: func() (float64, error) { return cbr.NewClient().GetRate("USD", time.Now()) },
	}

	updateRates()
	startCmdOnSchedule(updateRates)

	run()
}

func updateRates() {
	if err := forex.updateRate(); err != nil {
		log.Error(err)
	}
	log.Debugf(forex.pattern, forex.rate)

	if err := cbrf.updateRate(); err != nil {
		log.Error(err)
	}
	log.Debugf(cbrf.pattern, cbrf.rate)
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
				persist(update.Message.From)
				msg.Text = homeCmd
				msg.ReplyMarkup = homeKeyboard
			case "forex":
				msg.Text = forex.getRatef()
			case "cbrf":
				msg.Text = cbrf.getRatef()
			case "help":
				persist(update.Message.From)
				msg.Text = helpCmd
			default:
				msg.Text = unknownCmd
			}

			// Send the message.
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
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, forex.getRatef())
			case cbrTitle:
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, cbrf.getRatef())
			case "Help":
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, helpCmd)
			default:
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			}
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
