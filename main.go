package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Send any text message to the bot after the bot has been started

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithCallbackQueryDataHandler("asset", bot.MatchTypePrefix, callbackHandler),
		bot.WithCallbackQueryDataHandler("buyAndSell", bot.MatchTypePrefix, clickHandler),
		bot.WithCallbackQueryDataHandler("wallet", bot.MatchTypePrefix, clickHandler),
		bot.WithCallbackQueryDataHandler("setting", bot.MatchTypePrefix, clickHandler),
		bot.WithCallbackQueryDataHandler("analysis", bot.MatchTypePrefix, clickHandler),
		bot.WithCallbackQueryDataHandler("refresh", bot.MatchTypePrefix, clickHandler),
		bot.WithCallbackQueryDataHandler("referral", bot.MatchTypePrefix, clickHandler),
		bot.WithCallbackQueryDataHandler("language", bot.MatchTypePrefix, clickHandler),
	}

	//b, err := bot.New(os.Getenv("EXAMPLE_TELEGRAM_BOT_TOKEN"), opts...)
	b, err := bot.New("7869146013:AAGYL-txb0sGTDOUUQe5B0yDOZy-roTD0EY", opts...)

	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		panic(err)
	}

	b.Start(ctx)
}

func callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// answering callback query first to let Telegram know that we received the callback query,
	// and we're handling it. Otherwise, Telegram might retry sending the update repetitively
	// as it thinks the callback query doesn't reach to our application. learn more by
	// reading the footnote of the https://core.telegram.org/bots/api#callbackquery type.
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the button: " + update.CallbackQuery.Data,
	})
}
func clickHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// answering callback query first to let Telegram know that we received the callback query,
	// and we're handling it. Otherwise, Telegram might retry sending the update repetitively
	// as it thinks the callback query doesn't reach to our application. learn more by
	// reading the footnote of the https://core.telegram.org/bots/api#callbackquery type.
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You click the button: " + update.CallbackQuery.Data,
	})
}
func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "💰买/卖", CallbackData: "buyAndSell"},
				{Text: "💵资产", CallbackData: "asset"},
			}, {
				{Text: "💳我的钱包", CallbackData: "wallet"},
				{Text: "⚙️设置", CallbackData: "setting"},
			}, {
				{Text: "钱包分析", CallbackData: "analysis"},
				{Text: "刷新余额", CallbackData: "refresh"},
			}, {
				{Text: "🎁邀请返佣", CallbackData: "referral"},
				{Text: "🌎Language/语言", CallbackData: "language"},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "☀️欢迎使用DT bsc交易机器人☀️\n🚀当前交易钱包：\n💳钱包余额：0.000000 BNB",
		ReplyMarkup: kb,
	})
}
