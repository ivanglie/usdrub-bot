package main

import (
	"fmt"
	"os"
	"strings"
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

	fx, mx, cbrf                    *ex.Currency
	cash                            *cashex.Currency
	buyBranches, sellBranches       map[int][]string
	sellChatIdIndex, buyChatIdIndex map[int64]int
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
	scheduler.Debug, forex.Debug, moex.Debug, cbr.Debug, cashex.Debug = opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg

	tgbotapi.SetLogger(log)
	scheduler.SetLogger(log)
	forex.SetLogger(log)
	cbr.SetLogger(log)
	moex.SetLogger(log)
	cashex.SetLogger(log)

	mx = ex.New(func() (float64, error) { return moex.NewClient().GetRate(moex.USDRUB) })
	fx = ex.New(func() (float64, error) { return forex.NewClient().GetRate("USD", "RUB") })
	cbrf = ex.New(func() (float64, error) { return cbr.NewClient().GetRate("USD", time.Now()) })
	cash = cashex.New(cashex.Region)

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

		if update.Message != nil {
			if !update.Message.IsCommand() {
				continue
			}

			commandHandler(update)
		} else if update.CallbackQuery != nil {
			callbackQueryHandler(update)
		}
	}
}

func commandHandler(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	switch update.Message.Command() {
	case "start", "dashboard":
		err := storage.Persist(update.Message.From)
		if err != nil {
			log.Error(err)
		}
		msg.Text = fmt.Sprintf("*%s*\n%s %s\n%s %s\n%s %s\n*%s*\n%s\n%s",
			exPrefix, fx, fxSuffix, mx, mxSuffix, cbrf, cbrfSuffix, cashPrefix, cash, cashSuffix)
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
		err := storage.Persist(update.Message.From)
		if err != nil {
			log.Error(err)
		}
		msg.Text = helpCmd
	default:
		msg.Text = unknownCmd
	}

	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := bot.Send(msg); err != nil {
		log.Error(err)
	}
}

func callbackQueryHandler(update tgbotapi.Update) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := bot.Request(callback); err != nil {
		log.Error(err)
	}

	msg := tgbotapi.MessageConfig{}
	switch update.CallbackQuery.Data {
	case "Buy":
		buyBranches = cash.BuyBranches()

		chatId := update.CallbackQuery.Message.Chat.ID
		buyChatIdIndex = make(map[int64]int)
		buyChatIdIndex[chatId] = 0
		msg = tgbotapi.NewMessage(chatId, "*Buy cash*\n"+strings.Join(buyBranches[buyChatIdIndex[chatId]], "\n"))
		if len(buyBranches) > 1 {
			msg.ReplyMarkup = cashBuyMoreKeyboard
		}
	case "BuyMore":
		chatId := update.CallbackQuery.Message.Chat.ID
		if buyChatIdIndex[chatId] < len(buyBranches) {
			buyChatIdIndex[chatId] = buyChatIdIndex[chatId] + 1
		}

		if buyBranches[buyChatIdIndex[chatId]] != nil {
			msg = tgbotapi.NewMessage(chatId, "*Buy cash*\n"+strings.Join(buyBranches[buyChatIdIndex[chatId]], "\n"))
		}

		if buyChatIdIndex[chatId] != len(buyBranches)-1 {
			msg.ReplyMarkup = cashBuyMoreKeyboard
		}
	case "Sell":
		sellBranches = cash.SellBranches()

		chatId := update.CallbackQuery.Message.Chat.ID
		sellChatIdIndex = make(map[int64]int)
		sellChatIdIndex[chatId] = 0
		msg = tgbotapi.NewMessage(chatId, "*Sell cash*\n"+strings.Join(sellBranches[sellChatIdIndex[chatId]], "\n"))
		if len(sellBranches) > 1 {
			msg.ReplyMarkup = cashSellMoreKeyboard
		}
	case "SellMore":
		chatId := update.CallbackQuery.Message.Chat.ID
		if sellChatIdIndex[chatId] < len(sellBranches) {
			sellChatIdIndex[chatId] = sellChatIdIndex[chatId] + 1
		}

		if sellBranches[sellChatIdIndex[chatId]] != nil {
			msg = tgbotapi.NewMessage(chatId, "*Sell cash*\n"+strings.Join(sellBranches[sellChatIdIndex[chatId]], "\n"))
		}

		if sellChatIdIndex[chatId] != len(sellBranches)-1 {
			msg.ReplyMarkup = cashSellMoreKeyboard
		}
	case "Help":
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, helpCmd)
	default:
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	}
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := bot.Send(msg); err != nil {
		log.Error(err)
	}
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
