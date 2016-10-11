package main

import (
	"fmt"
	"html"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func search(msg string) (string, error) {
	msg = strings.Trim(msg, " ")
	uri := fmt.Sprintf("%s/search?type=content&q=%s", cfg.Zhihu.Host, url.QueryEscape(msg))
	doc, err := goquery.NewDocument(uri)
	if err != nil {
		log.Println(err)
		return "", err
	}

	msg = ""
	doc.Find("ul.list li").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".title").Text()
		smy := s.Find(".content .summary")
		smy.Find("a.toggle-expand").Remove()
		summary := smy.Text()
		// content := s.Find(".visible-expanded .content").Text()

		questionLink, _ := s.Find("a").Attr("href")
		answerLink, _ := s.Find(".entry-body .entry-content").Attr("data-entry-url")
		if title == "" {
			return
		}

		msg = fmt.Sprintf(`%s<a href="%s/%s">%s</a><br>%s <a href="%s/%s">...显示全部</a><br><br>`,
			msg, cfg.Zhihu.Host, questionLink, title, html.EscapeString(summary), cfg.Zhihu.Host, answerLink)
	})

	msg = format(msg)
	return msg, nil
}

var (
	Warp = `
	`
	ReplaceHTML = map[string]string{
		"<br>":       Warp,
		"&lt;br&gt;": Warp,
	}
)

func format(msg string) string {
	for k, v := range ReplaceHTML {
		msg = strings.Replace(msg, k, v, -1)
	}
	return msg
}
