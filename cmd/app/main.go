package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/go-telegram/ui/paginator"
	"github.com/ivanglie/go-br-client"
	"github.com/ivanglie/go-cbr-client"
	"github.com/ivanglie/go-coingate-client"
	"github.com/ivanglie/go-moex-client"
	"github.com/ivanglie/usdrub-bot/internal/exrate"
	"github.com/ivanglie/usdrub-bot/internal/utils"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

const (
	helpCmd    = "Just use /forex, /moex, /cbrf, /cash and /dashboard command."
	exPrefix   = "1 US Dollar equals"
	cashPrefix = "Exchange rates of cash"
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

	forexRate,
	moexRate,
	cbrfRate *exrate.Rate
	cashRate *exrate.CashRate

	forexRateCh,
	moexRateCh,
	cbrfRateCh chan *exrate.Rate
	cashRateCh chan *exrate.CashRate
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
	coingate.Debug, moex.Debug, cbr.Debug, br.Debug, utils.Debug = opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg, opts.Dbg

	forexRateCh = make(chan *exrate.Rate)
	moexRateCh = make(chan *exrate.Rate)
	cbrfRateCh = make(chan *exrate.Rate)
	cashRateCh = make(chan *exrate.CashRate)

	updateRates()
	if err := utils.StartCmdOnSchedule(updateRates); err != nil {
		log.Panic(err)
	}

	b, err := bot.New(opts.BotToken, bot.WithDebug())
	if err != nil {
		log.Panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/forex", bot.MatchTypeExact, forexHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/moex", bot.MatchTypeExact, moexHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cbrf", bot.MatchTypeExact, cbrfHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cash", bot.MatchTypeExact, cashHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, dashboardHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/dashboard", bot.MatchTypeExact, dashboardHandler)

	ctx := context.TODO()
	b.Start(ctx)
}

func updateRates() {
	t := time.Now()

	go func() {
		forexRateCh <- exrate.UpdateRate(func() (float64, error) { return coingate.NewClient().GetRate("USD", "RUB") })
	}()

	go func() {
		moexRateCh <- exrate.UpdateRate(func() (float64, error) { return moex.NewClient().GetRate(moex.USDRUB) })
	}()

	go func() {
		cbrfRateCh <- exrate.UpdateRate(func() (float64, error) { return cbr.NewClient().GetRate("USD", time.Now()) })
	}()

	go func() {
		cashRateCh <- exrate.UpdateCashRate(func() (*br.Rates, error) { return br.NewClient().Rates(br.USD, br.Moscow) })
	}()

	cashRate = <-cashRateCh
	cbrfRate = <-cbrfRateCh
	moexRate = <-moexRateCh
	forexRate = <-forexRateCh

	log.Debugln("Elapsed time:", time.Since(t))
}

func forexHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintln(exPrefix, forexRate, fxSuffix),
	})
}

func moexHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintln(exPrefix, moexRate, mxSuffix),
	})
}

func cbrfHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintln(exPrefix, cbrfRate, cbrfSuffix),
	})
}

func cashHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b).
		Row().
		Button("Buy cash", []byte("buy"), onBuy).
		Button("Sell cash", []byte("sell"), onSell).
		Button("Help", []byte("help"), onHelp)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        fmt.Sprintf("%s\n%s\n%s", cashPrefix, cashRate, cashSuffix),
		ReplyMarkup: kb,
	})
}

func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpCmd,
	})
}

func dashboardHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b, inline.NoDeleteAfterClick()).
		Row().
		Button("Buy cash", []byte("buy"), onBuy).
		Button("Sell cash", []byte("sell"), onSell).
		Button("Help", []byte("help"), onHelp)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf("<b>%s</b>\n%s %s\n%s %s\n%s %s\n<b>%s</b>\n%s\n%s",
			exPrefix, forexRate, fxSuffix, moexRate, mxSuffix, cbrfRate, cbrfSuffix, cashPrefix, cashRate, cashSuffix),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: kb,
	})

	if err := utils.Persist(update.Message.From); err != nil {
		log.Error(err)
	}
}

func onBuy(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	bb := cashRate.BuyBranches()
	s := []string{}
	for i, v := range bb {
		i++
		v = fmt.Sprintf("*%d* %s", i, bot.EscapeMarkdownUnescaped(v))
		s = append(s, v)
	}

	opts := []paginator.Option{
		paginator.PerPage(5),
		paginator.WithCloseButton("Close"),
	}

	log.Debugln(s)
	p := paginator.New(append([]string{"*Buy cash*"}, s...), opts...)

	p.Show(ctx, b, strconv.FormatInt(mes.Chat.ID, 10))
}

func onSell(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	sb := cashRate.SellBranches()
	s := []string{}
	for i, v := range sb {
		i++
		v := fmt.Sprintf("*%d* %s", i, bot.EscapeMarkdownUnescaped(v))
		s = append(s, v)
	}

	opts := []paginator.Option{
		paginator.PerPage(5),
		paginator.WithCloseButton("Close"),
	}

	log.Debugln(s)
	p := paginator.New(append([]string{"*Sell cash*"}, s...), opts...)

	p.Show(ctx, b, strconv.FormatInt(mes.Chat.ID, 10))
}

func onHelp(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Chat.ID,
		Text:   helpCmd,
	})
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
}
