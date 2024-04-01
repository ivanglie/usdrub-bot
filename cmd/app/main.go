package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ivanglie/usdrub-bot/internal/cexrate"
	"github.com/ivanglie/usdrub-bot/internal/exrate"
	"github.com/ivanglie/usdrub-bot/internal/logger"
	"github.com/ivanglie/usdrub-bot/internal/scheduler"
	"github.com/ivanglie/usdrub-bot/pkg/go-br-client"
	"github.com/ivanglie/usdrub-bot/pkg/go-cbr-client"
	"github.com/ivanglie/usdrub-bot/pkg/go-coingate-client"
	"github.com/ivanglie/usdrub-bot/pkg/go-moex-client"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

const (
	helpCmd = "Just use /forex, /moex, /cbrf, /cash and /dashboard command."
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
	coingate.Debug, moex.Debug, cbr.Debug, br.Debug, logger.Debug = opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg

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
				forex(bot, update)
			case "moex":
				moexHandler(bot, update)
			case "cbrf":
				cbrf(bot, update)
			case "cash":
				cash(bot, update)
			case "help":
				help(bot, update)
			case "start":
				start(bot, update)
			case "dashboard":
				dashboard(bot, update)
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

func forex(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Forex request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintln(exrate.Prefix, exrate.Get().Value(exrate.Forex)),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)

	bot.Send(msg)
}

func moexHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Moex request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintln(exrate.Prefix, exrate.Get().Value(exrate.MOEX)),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)

	bot.Send(msg)
}

func cbrf(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Cbrf request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintln(exrate.Prefix, exrate.Get().Value(exrate.CBRF)),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)

	bot.Send(msg)
}

func cash(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Cash request from %s", update.Message.From)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("<b>%s</b>\n%s\n%s", cexrate.Prefix, cexrate.Get().String(), cexrate.Suffix),
	)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = getReplyMessageID(update.Message)
	msg.ReplyMarkup = &kb

	bot.Send(msg)
}

func help(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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

	dashboard(bot, update)
}

func dashboard(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	log.Infof("Dashboard request from %s", update.Message.From)

	t := fmt.Sprintf("<b>%s</b>\n%s<b>%s</b>\n%s\n%s",
		exrate.Prefix, exrate.Get().String(),
		cexrate.Prefix, cexrate.Get().String(), cexrate.Suffix)

	if len(cexrate.Get().BuyBranches()) == 0 || len(cexrate.Get().SellBranches()) == 0 {
		t = fmt.Sprintf("<b>%s</b>\n%s", exrate.Prefix, exrate.Get().String())
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

	bb := cexrate.Get().BuyBranches()
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

	sb := cexrate.Get().SellBranches()
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
	br.SetLogger(log)
	logger.SetLogger(log)
}
