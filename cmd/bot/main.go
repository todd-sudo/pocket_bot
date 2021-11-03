package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/dev2033/go_tg_bot/pkg/repository"
	"github.com/dev2033/go_tg_bot/pkg/repository/boltdb"
	"github.com/dev2033/go_tg_bot/pkg/server"
	"github.com/dev2033/go_tg_bot/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("2008810566:AAEL5RS-GpbSil3JzjsLDBuYCccWjuIBsnk")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	pocketClient, err := pocket.NewClient("99472-07014b2b6f55884626dbf680")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost:8000/")

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https://t.me/pocket_golang_tg_bot")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
