package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

func botRun() error {
	var err error
	bot, err = tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = cfg.Bot.Debug

	log.Println("Authorized on account:", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert(fmt.Sprintf("%s%s/%s", cfg.HTTP.Host, cfg.HTTP.Port, cfg.Bot.Token), cfg.HTTP.PublicKey))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook(fmt.Sprintf("/%s", bot.Token))
	go func() {
		if err := http.ListenAndServeTLS(cfg.HTTP.Port, cfg.HTTP.PublicKey, cfg.HTTP.PrivateKey, nil); err != nil {
			log.Fatal(err)
		}
	}()

	for update := range updates {
		msgRouter(update)
	}
	return nil
}

func msgRouter(update tgbotapi.Update) error {
	switch {
	case update.InlineQuery != nil:
		return isInline(update)
	case update.Message != nil && update.Message.IsCommand():
		return isCommand(update)
	case update.Message != nil && (update.Message.Chat.IsPrivate() || bot.IsMessageToMe(*update.Message)):
		return isSearch(update)
	}
	return nil
}

func isCommand(update tgbotapi.Update) error {
	switch update.Message.Command() {
	case "s":
		return isSearch(update)
	case "daily":
		return isDaily(update)
	default:
		return sendMsg(update, HelpMsg)
	}
	return nil
}

func isSearch(update tgbotapi.Update) error {
	var msg string
	if update.Message.IsCommand() {
		msg = update.Message.CommandArguments()
	} else {
		msg = update.Message.Text
	}
	msg = strings.Trim(msg, " ")
	if msg == "" {
		return sendMsg(update, HelpMsg)
	}

	results, err := search(update.Message.Text, cfg.Zhihu.SearchResultNum)
	if err != nil {
		return err
	}

	msg = ""
	for _, result := range results {
		msg = fmt.Sprintf(`%s<a href="%s">%s</a><br>%s <a href="%s">...显示全部</a><br><br>`,
			msg, result.QuestionLink, result.Title, html.EscapeString(result.Summary), result.AnswerLink)
	}
	msg = format(msg)
	return sendMsg(update, msg)
}

func isInline(update tgbotapi.Update) error {
	msg := update.InlineQuery.Query
	results, err := search(msg, cfg.Zhihu.InlineResultNum)
	if err != nil {
		return err
	}
	var answers []interface{}
	for _, result := range results {
		content := html.EscapeString(result.Content)
		if len(content) > 2000 {
			content = Substr(content, 2000)
		}
		msg = fmt.Sprintf(`<a href="%s">%s</a><br>%s <a href="%s">...显示全部</a><br><br>`,
			result.QuestionLink, result.Title, content, result.AnswerLink)
		msg = format(msg)
		answer := tgbotapi.NewInlineQueryResultArticleHTML(result.ID, result.Title, msg)
		answer.Description = html.EscapeString(result.Summary)
		inputTextMessageContent := answer.InputMessageContent.(tgbotapi.InputTextMessageContent)
		inputTextMessageContent.DisableWebPagePreview = true
		answer.InputMessageContent = inputTextMessageContent
		answers = append(answers, &answer)
	}
	return answerInlineQuery(update, answers)
}

func isDaily(update tgbotapi.Update) error {
	txt, err := daily()
	if err != nil {
		return err
	}
	return sendMsg(update, txt)
}

func sendMsg(update tgbotapi.Update, txt string) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, txt)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	if resp, err := bot.Send(msg); err != nil {
		log.Println("bot.Send:", err, resp)
		return err
	}
	return nil
}

func answerInlineQuery(update tgbotapi.Update, results []interface{}) error {
	answer := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       results,
	}
	if resp, err := bot.AnswerInlineQuery(answer); err != nil {
		log.Println("bot.answerInlineQuery:", err, resp)
		return err
	}
	return nil
}
