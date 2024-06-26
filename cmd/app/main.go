package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanglie/usdrub-bot/internal/cash"
	"github.com/ivanglie/usdrub-bot/internal/crypto"
	"github.com/ivanglie/usdrub-bot/internal/exchange"
	"github.com/ivanglie/usdrub-bot/internal/logger"
	"github.com/ivanglie/usdrub-bot/internal/scheduler"
	"github.com/ivanglie/usdrub-bot/pkg/bankiru-go"
	"github.com/ivanglie/usdrub-bot/pkg/bestchange-go"
	"github.com/ivanglie/usdrub-bot/pkg/cbr-go"
	"github.com/ivanglie/usdrub-bot/pkg/coingate-go"
	"github.com/ivanglie/usdrub-bot/pkg/moex-go"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

const (
	helpCmd = "Just use /forex, /moex, /cbrf, /cash, /crypto and /dashboard command."
)

var (
	log *logrus.Logger

	opts struct {
		Dbg      bool   `long:"dbg" env:"DEBUG" description:"Debug mode"`
		BotToken string `long:"bottoken" env:"BOT_TOKEN" description:"Telegram API Token"`
		CronSpec string `long:"cronspec" env:"CRON_SPEC" description:"Cron spec"`
	}

	kb = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Buy cash", "Buy"),
			tgbotapi.NewInlineKeyboardButtonData("Sell cash", "Sell"),
			tgbotapi.NewInlineKeyboardButtonData("Help", "Help"),
		),
	)

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
	coingate.Debug, moex.Debug, cbr.Debug, bankiru.Debug, bestchange.Debug, logger.Debug = opts.Dbg, opts.Dbg, opts.Dbg,
		opts.Dbg, opts.Dbg, opts.Dbg

	updateRates := func() {
		t := time.Now()

		type RateInterface interface {
			Update()
		}

		rates := []RateInterface{exchange.Get(), cash.Get(), crypto.Get()}

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

	if err := scheduler.StartCmdOnSchedule(updateRates); err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(opts.BotToken)
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

			switch update.Message.Command() {
			case "forex":
				forexHandler(bot, update)
			case "moex":
				moexHandler(bot, update)
			case "cbrf":
				cbrfHandler(bot, update)
			case "cash":
				cashHandler(bot, update)
			case "crypto":
				cryptoHandler(bot, update)
			case "help":
				helpHandler(bot, update)
			case "start":
				start(bot, update)
			case "dashboard":
				dashboardHandler(bot, update)
			default:
				log.Warnf("Unknown command %q", update.Message.Command())
			}
		}

		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "Buy":
				onBuy(bot, update.CallbackQuery)
			case "Sell":
				onSell(bot, update.CallbackQuery)
			case "Help":
				onHelp(bot, update.CallbackQuery)
			default:
				log.Warnf("Unknown callback %q", update.CallbackQuery.Data)
			}
		}
	}
}

func forexHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Forex request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintln(exchange.Prefix, exchange.Get().Value(exchange.Forex)),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)

	bot.Send(msg)
}

func moexHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Moex request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintln(exchange.Prefix, exchange.Get().Value(exchange.MOEX)),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)

	bot.Send(msg)
}

func cbrfHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Cbrf request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintln(exchange.Prefix, exchange.Get().Value(exchange.CBRF)),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)

	bot.Send(msg)
}

func cashHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Cash request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("<b>%s</b>\n%s\n%s", cash.Prefix, cash.Get().String(), cash.Suffix),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)
	msg.ReplyMarkup = &kb

	bot.Send(msg)
}

func cryptoHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Crypto request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("<b>%s</b>\n%s", crypto.Prefix, crypto.Get().String()),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)
	msg.ReplyMarkup = &kb

	bot.Send(msg)
}

func helpHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Help request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		helpCmd,
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)

	bot.Send(msg)
}

func start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Start request from %s", update.Message.From)

	dashboardHandler(bot, update)
}

func dashboardHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Dashboard request from %s", update.Message.From)

	t := fmt.Sprintf("<b>%s</b>\n%s<b>%s</b>\n%s\n<b>%s</b>\n%s\n%s",
		exchange.Prefix, exchange.Get().String(),
		crypto.Prefix, crypto.Get().String(),
		cash.Prefix, cash.Get().String(), cash.Suffix)

	if len(cash.Get().BuyBranches()) == 0 || len(cash.Get().SellBranches()) == 0 {
		t = fmt.Sprintf("<b>%s</b>\n%s", exchange.Prefix, exchange.Get().String())
		log.Warn("No branches")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, t)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)
	msg.ReplyMarkup = &kb

	bot.Send(msg)
}

func onBuy(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	log.Infof("OnBuy request from %s", cq.From)

	bb := cash.Get().BuyBranches()
	if len(bb) == 0 {
		log.Warn("No buy branches")
		return
	}

	s := []string{}
	for i, v := range bb {
		i++
		v = fmt.Sprintf("<b>%d</b> %s", i, v)
		s = append(s, v)
	}

	log.Debugln(s)

	msg := tgbotapi.NewMessage(
		cq.Message.Chat.ID,
		strings.Join(append([]string{"<b>Buy cash</b>"}, s...), "\n"),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(cq.Message)

	bot.Send(msg)
}

func onSell(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	log.Infof("OnSell request from %s", cq.From)

	sb := cash.Get().SellBranches()
	if len(sb) == 0 {
		log.Warn("No sell branches")
		return
	}

	s := []string{}
	for i, v := range sb {
		i++
		v := fmt.Sprintf("<b>%d</b> %s", i, v)
		s = append(s, v)
	}

	log.Debugln(s)

	msg := tgbotapi.NewMessage(
		cq.Message.Chat.ID,
		strings.Join(append([]string{"<b>Sell cash</b>"}, s...), "\n"),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(cq.Message)

	bot.Send(msg)
}

func onHelp(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	log.Infof("OnHelp request from %s", cq.From)

	msg := tgbotapi.NewMessage(
		cq.Message.Chat.ID,
		helpCmd,
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(cq.Message)

	bot.Send(msg)
}

// getReplyMessageID returns message to reply to.
func getReplyMessageID(message *tgbotapi.Message) int {
	if message.Chat.Type != "private" {
		return message.MessageID
	}

	return 0
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

	log.SetLevel(logrus.InfoLevel)
}

func setLogger(log *logrus.Logger) {
	coingate.SetLogger(log)
	cbr.SetLogger(log)
	moex.SetLogger(log)
	bankiru.SetLogger(log)
	bestchange.SetLogger(log)
	logger.SetLogger(log)
}
