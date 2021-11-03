package telegram

import (
	"log"

	"github.com/dev2033/go_tg_bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	redirectURL     string
}

// NewBot constructor
func NewBot(
	bot *tgbotapi.BotAPI,
	pocketClient *pocket.Client,
	tr repository.TokenRepository,
	redirectURL string,
) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, redirectURL: redirectURL, tokenRepository: tr}
}

// Start bot
func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.handlerUpdates(updates)

	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handlerUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		b.handleMessage(update.Message)
	}
}
