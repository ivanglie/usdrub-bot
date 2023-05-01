package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"github.com/go-telegram/ui/paginator"
	"github.com/ivanglie/usdrub-bot/internal/cexrate"
	"github.com/ivanglie/usdrub-bot/internal/exrate"
	"github.com/ivanglie/usdrub-bot/internal/utils"
)

const (
	helpCmd = "Just use /forex, /moex, /cbrf, /cash and /dashboard command."
)

func CreateAndStart(token string) {
	opts := []bot.Option{}
	if Debug {
		opts = append(opts, bot.WithDebug())
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Printf("[ERROR] %v", err)
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

func forexHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             fmt.Sprintln(exrate.Prefix, exrate.Get().Value(exrate.Forex)),
		ReplyToMessageID: getReplyMessageID(update.Message),
	})
}

func moexHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             fmt.Sprintln(exrate.Prefix, exrate.Get().Value(exrate.MOEX)),
		ReplyToMessageID: getReplyMessageID(update.Message),
	})
}

func cbrfHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             fmt.Sprintln(exrate.Prefix, exrate.Get().Value(exrate.CBRF)),
		ReplyToMessageID: getReplyMessageID(update.Message),
	})
}

func cashHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:           update.Message.Chat.ID,
		Text:             helpCmd,
		ReplyToMessageID: getReplyMessageID(update.Message),
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
		Text: fmt.Sprintf("*%s*\n%s*%s*\n%s\n%s",
			exrate.Prefix, bot.EscapeMarkdownUnescaped(exrate.Get().String()),
			cexrate.Prefix, bot.EscapeMarkdownUnescaped(cexrate.Get().String()), bot.EscapeMarkdownUnescaped(cexrate.Suffix)),
		ParseMode:        models.ParseModeMarkdown,
		ReplyMarkup:      kb,
		ReplyToMessageID: getReplyMessageID(update.Message),
	})

	if err := utils.Persist(update.Message.From); err != nil {
		log.Printf("[ERROR] %v", err)
	}
}

func onBuy(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
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

	log.Printf("[DEBUG] %s", s)
	p := paginator.New(append([]string{"*Buy cash*"}, s...), opts...)

	p.Show(ctx, b, strconv.FormatInt(mes.Chat.ID, 10))
}

func onSell(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
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

	log.Printf("[DEBUG] %s", s)
	p := paginator.New(append([]string{"*Sell cash*"}, s...), opts...)

	p.Show(ctx, b, strconv.FormatInt(mes.Chat.ID, 10))
}

func onHelp(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
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
