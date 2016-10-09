package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

func botRun() error {
	var err error
	bot, err = tgbotapi.NewBotAPI(cfg.Bot.Token)
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
	if update.Message == nil {
		log.Println("message is nil", update)
		return nil
	}

	switch {
	case update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message):
		sendMsg(update)
	}
	return nil
}

func sendMsg(update tgbotapi.Update) error {
	text, err := search(update.Message.Text)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	if _, err = bot.Send(msg); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
