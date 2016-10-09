package main

import (
	"log"
	"net/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	token = "266361856:AAGNWqZLAw2DVKUtEeTcHT_mZS2t1kkIV00"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://bot.l2x.me:8443/"+bot.Token, "./public.pem"))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go func() {
		if err := http.ListenAndServeTLS(":8443", "public.pem", "private.key", nil); err != nil {
			panic(err)
		}
	}()

	for update := range updates {
		log.Printf("%+v\n", update)
	}
}
