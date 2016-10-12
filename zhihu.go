package main

import (
	"fmt"
	"html"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	ID           string
	Title        string
	Summary      string
	Content      string
	QuestionLink string
	AnswerLink   string
}

func search(msg string, limit int) ([]SearchResult, error) {
	uri := fmt.Sprintf("%s/search?type=content&q=%s", cfg.Zhihu.Host, url.QueryEscape(msg))
	doc, err := goquery.NewDocument(uri)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var results []SearchResult
	doc.Find("ul.list li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i >= limit {
			return false
		}

		title := s.Find(".title").Text()
		smy := s.Find(".content .summary")
		smy.Find("a.toggle-expand").Remove()
		summary := format(smy.Text())
		content := format(s.Find(".visible-expanded .content").Text())

		questionLink, _ := s.Find("a").Attr("href")
		answerLink, _ := s.Find(".entry-body .entry-content").Attr("data-entry-url")
		tmp := strings.Split(answerLink, "/")
		id := tmp[len(tmp)-1]

		if !strings.HasPrefix(questionLink, "http") {
			questionLink = fmt.Sprintf("%s/%s", cfg.Zhihu.Host, strings.TrimLeft(questionLink, "/"))
		}
		if !strings.HasPrefix(answerLink, "http") {
			answerLink = fmt.Sprintf("%s/%s", cfg.Zhihu.Host, strings.TrimLeft(answerLink, "/"))
		}
		if title == "" {
			return true
		}

		result := SearchResult{
			ID:           id,
			Title:        title,
			Summary:      html.EscapeString(summary),
			Content:      html.EscapeString(content),
			QuestionLink: questionLink,
			AnswerLink:   answerLink,
		}
		results = append(results, result)

		return true
	})

	return results, nil
}

func daily() (string, error) {
	return "", nil
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

	re := regexp.MustCompile("<(/?)[a-zA-Z+]>")
	msg = re.ReplaceAllString(msg, "")
	return msg
}

func Substr(s string, size int) string {
	if len(s) < size {
		return s
	}
	var k, i int
	chars := []rune(s)
	for _, c := range chars {
		if c > 254 {
			i += 3
		} else {
			i++
		}
		if i > size {
			break
		}
		k++
	}
	return string(chars[0:k])
}
