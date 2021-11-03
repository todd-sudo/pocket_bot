package telegram

import (
	"context"
	"fmt"

	"github.com/dev2033/go_tg_bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// initAuthorizationProcess запускает процесс авторизации
func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.createAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}

	msgText := fmt.Sprintf(replyStartTempate, authLink)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) createAuthorizationLink(chatID int64) (string, error) {
	redirectUrl := b.generateRedirectLink(chatID)
	token, err := b.pocketClient.GetRequestToken(context.Background(), b.redirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, token, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(token, redirectUrl)
}

// getAccessToken получает access token
func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.tokenRepository.Get(chatID, repository.AccessTokens)
}

// generateRedirectLink генерирует ссылку для редиректа
func (b *Bot) generateRedirectLink(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
