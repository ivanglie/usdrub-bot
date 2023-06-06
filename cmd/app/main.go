package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/go-telegram/ui/paginator"
	"github.com/ivanglie/go-br-client"
	"github.com/ivanglie/go-cbr-client"
	"github.com/ivanglie/go-coingate-client"
	"github.com/ivanglie/go-moex-client"
	"github.com/ivanglie/usdrub-bot/internal/cexrate"
	"github.com/ivanglie/usdrub-bot/internal/exrate"
	"github.com/ivanglie/usdrub-bot/internal/logger"
	"github.com/ivanglie/usdrub-bot/internal/scheduler"

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

	bOpts := []bot.Option{}
	if opts.Dbg {
		bOpts = append(bOpts, bot.WithDebug())
	}

	b, err := bot.New(opts.BotToken, bOpts...)
	if err != nil {
		log.Panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/forex", bot.MatchTypePrefix, forexHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/moex", bot.MatchTypePrefix, moexHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cbrf", bot.MatchTypePrefix, cbrfHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cash", bot.MatchTypePrefix, cashHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypePrefix, helpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/dashboard", bot.MatchTypePrefix, dashboardHandler)

	ctx := context.TODO()
	b.Start(ctx)
}

func forexHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logInfo("Request forex", update.Message)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             fmt.Sprintln(exrate.Prefix, exrate.Get().Value(exrate.Forex)),
		ReplyToMessageID: getReplyMessageID(update.Message),
	})
}

func moexHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logInfo("Request moex", update.Message)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             fmt.Sprintln(exrate.Prefix, exrate.Get().Value(exrate.MOEX)),
		ReplyToMessageID: getReplyMessageID(update.Message),
	})
}

func cbrfHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logInfo("Request cbrf", update.Message)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             fmt.Sprintln(exrate.Prefix, exrate.Get().Value(exrate.CBRF)),
		ReplyToMessageID: getReplyMessageID(update.Message),
	})
}

func cashHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logInfo("Request cash", update.Message)

	kb := inline.New(b).
		Row().
		Button("Buy cash", []byte("buy"), onBuy).
		Button("Sell cash", []byte("sell"), onSell).
		Button("Help", []byte("help"), onHelp)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             fmt.Sprintf("%s\n%s\n%s", cexrate.Prefix, cexrate.Get().String(), cexrate.Suffix),
		ReplyMarkup:      kb,
		ReplyToMessageID: getReplyMessageID(update.Message),
	})
}

func helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logInfo("Request help", update.Message)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             helpCmd,
		ReplyToMessageID: getReplyMessageID(update.Message),
	})
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logInfo("Request start", update.Message)

	dashboardHandler(ctx, b, update)
}

func dashboardHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	logInfo("Request dashboard", update.Message)

	kb := inline.New(b, inline.NoDeleteAfterClick()).
		Row().
		Button("Buy cash", []byte("buy"), onBuy).
		Button("Sell cash", []byte("sell"), onSell).
		Button("Help", []byte("help"), onHelp)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf("<b>%s</b>\n%s<b>%s</b>\n%s\n%s",
			exrate.Prefix, exrate.Get().String(),
			cexrate.Prefix, cexrate.Get().String(), cexrate.Suffix),
		ParseMode:        models.ParseModeHTML,
		ReplyMarkup:      kb,
		ReplyToMessageID: getReplyMessageID(update.Message),
	})
}

func onBuy(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	logInfo("Request onBuy", mes)

	bb := cexrate.Get().BuyBranches()
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
	logInfo("Request onSell", mes)

	sb := cexrate.Get().SellBranches()
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
	logInfo("Request onHelp", mes)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           mes.Chat.ID,
		Text:             helpCmd,
		ReplyToMessageID: getReplyMessageID(mes),
	})
}

// getReplyMessageID returns message to reply to.
func getReplyMessageID(message *models.Message) int {
	if message.Chat.Type != "private" {
		return message.ID
	}

	return 0
}

// logInfo logs info about message.
func logInfo(info string, mes *models.Message) {
	log.Infoln(info)

	user := mes.From
	if user == nil {
		return
	}

	log.Infof("User"+
		" (id: %d,"+
		" is_bot: %t,"+
		" first_name: %s,"+
		" last_name: %s,"+
		" username: %s,"+
		" language_code: %s,"+
		" is_premium: %t,"+
		" added_to_attachment_menu: %t,"+
		" can_join_groups: %t,"+
		" can_read_all_group_messages: %t,"+
		" support_inline_queries: %t)",
		user.ID, user.IsBot, user.FirstName, user.LastName, user.Username, user.LanguageCode, user.IsPremium,
		user.AddedToAttachmentMenu, user.CanJoinGroups, user.CanReadAllGroupMessages, user.SupportInlineQueries)

	chat := &mes.Chat
	if chat == nil {
		return
	}

	log.Infof("Chat"+
		" (id: %d,"+
		" type: %s,"+
		" title: %s,"+
		" bio: %s,"+
		" description: %s,"+
		" invite_link: %s,"+
		" username: %s,"+
		" first_name: %s,"+
		" last_name: %s)",
		chat.ID, chat.Type, chat.Title, chat.Bio, chat.Description, chat.InviteLink, chat.Username, chat.FirstName,
		chat.LastName)
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
