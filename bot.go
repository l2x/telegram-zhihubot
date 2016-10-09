package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func botRun() error {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return err
	}

	bot.Debug = cfg.Bot.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert(fmt.Sprintf("%s%s/%s", cfg.HTTP.Host, cfg.HTTP.Port, cfg.Bot.Token), cfg.HTTP.PublicKey))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook(fmt.Sprintf("/%s", bot.Token))
	go func() {
		if err := http.ListenAndServeTLS(cfg.HTTP.Port, cfg.HTTP.PublicKey, cfg.HTTP.PrivateKey, nil); err != nil {
			panic(err)
		}
	}()

	for update := range updates {
		msgRouter(update)
	}
	return nil
}

func msgRouter(update tgbotapi.Update) error {
	fmt.Println(update)
	switch update.Message {
	}
	return nil
}
