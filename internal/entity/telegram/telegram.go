package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Telegram struct {
	bot    *tgbotapi.BotAPI
	logger *zap.Logger
}

func NewTelegram(token string, logger *zap.Logger) (t *Telegram, err error) {
	t = &Telegram{bot: nil, logger: logger}
	if t.bot, err = tgbotapi.NewBotAPI(token); err != nil {
		return nil, err
	}

	t.logger.Info("authorized telegram service", zap.String("account", t.bot.Self.UserName))
	return t, nil
}

func (t *Telegram) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)
	defer t.bot.StopReceivingUpdates()

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			t.logger.Info(fmt.Sprintf("[%s] %s (id: %v, chatId: %v)", update.Message.From.UserName, update.Message.Text, update.Message.From.ID, update.Message.Chat.ID))

			// for now, just mirror incoming messages to make sure that the bot works
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			_, _ = t.bot.Send(msg)

		case <-ctx.Done():
			t.logger.Info("telegram shutdown")
			return
		}
	}
}

func (t *Telegram) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.bot.Send(msg)
	return err
}
