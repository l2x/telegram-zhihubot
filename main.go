package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	msg := "如何装逼"
	msg = search(msg)

	fmt.Println(msg)
}

var (
	ZhihuHost = "https://www.zhihu.com"
)

func search(msg string) string {
	uri := fmt.Sprintf("%s/search?type=content&sort=upvote&q=%s", ZhihuHost, url.QueryEscape(msg))
	doc, err := goquery.NewDocument(uri)
	if err != nil {
		log.Fatal(err)
	}
	msg = ""
	doc.Find("ul.list li").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".title").Text()
		smy := s.Find(".content .summary")
		smy.Find("a.toggle-expand").Remove()
		summary := smy.Text()
		content := s.Find(".visible-expanded .content").Text()

		questionLink, _ := s.Find("a").Attr("href")
		answerLink, _ := s.Find(".entry-body .entry-content").Attr("data-entry-url")

		fmt.Println(title)
		fmt.Println(questionLink)

		fmt.Println(summary)
		fmt.Println(answerLink)

		fmt.Println(content)

		msg = fmt.Sprintf(`%s<a href="%s/%s">%s</a><br>%s <a href="%s/%s">...显示全部</a><br><br>`, msg, ZhihuHost, questionLink, title, summary, ZhihuHost, answerLink)
	})

	return format(msg)
}

var (
	Warp = `
	`
	ReplaceHTML = map[string]string{
		"<p>":  "",
		"</p>": "",
		"<br>": Warp,
	}
)

func format(msg string) string {
	for k, v := range ReplaceHTML {
		msg = strings.Replace(msg, k, v, -1)
	}
	return msg
}
