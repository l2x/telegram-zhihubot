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
	case update.Message.IsCommand():
		return isCommand(update)
	case update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message):
		return isMessage(update)
	}
	return nil
}

func isCommand(update tgbotapi.Update) error {
	switch update.Message.Command() {
	case "s":
		txt, err := search(update.Message.CommandArguments())
		if err != nil {
			return err
		}
		return sendMsg(update, txt)
	}
	return nil
}

func isMessage(update tgbotapi.Update) error {
	txt, err := search(update.Message.Text)
	if err != nil {
		return err
	}
	return sendMsg(update, txt)
}

func sendMsg(update tgbotapi.Update, txt string) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
