package scraper

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var skipTexts = map[string]bool{
	"": true,
	"Source: Webnovel.com, updated by novlove.com": true,
}

type Chapter struct {
	Text string
	Prev *string
	Next *string
}

func GetNovel(url *string) Chapter {
	resp, err := http.Get(*url)
	if err != nil {
		return Chapter{}
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	var sb strings.Builder
	parseAndAppend(doc, "#chr-content h4", &sb, nil)
	parseAndAppend(doc, "#chr-content p", &sb, skipTexts)

	link_next := doc.Find("#next_chap").AttrOr("href", "")
	link_prev := doc.Find("#prev_chap").AttrOr("href", "")

	return Chapter{
		Text: strings.TrimSpace(sb.String()),
		Next: &link_next,
		Prev: &link_prev,
	}

}

func parseAndAppend(doc *goquery.Document, selector string, sb *strings.Builder, skipTexts map[string]bool) {
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if skipTexts != nil {
			if _, skip := skipTexts[text]; skip {
				return
			}
		}

		var buf bytes.Buffer
		for _, node := range s.Nodes {
			if err := html.Render(&buf, node); err != nil {
				sb.WriteString(text)
			} else {
				sb.Write(buf.Bytes())
			}
		}
		sb.WriteString("\n")
	})
}
