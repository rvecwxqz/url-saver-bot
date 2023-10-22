package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"unicode"
)

type parser struct {
	client *http.Client
}

func NewParser() parser {
	return parser{
		client: &http.Client{},
	}
}

func (p parser) parse(url string) (string, error) {
	response, err := p.client.Get(url)
	if err != nil || response.StatusCode >= 400 {
		return "", fmt.Errorf("error parsing URL %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing document %w", err)
	}
	response.Body.Close()

	texts := make([]string, 0, 100)
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		texts = append(texts, s.Text())
	})
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		texts = append(texts, s.Text())
	})

	if len(texts) == 0 {
		return "", NewNoDataError()
	}

	outTexts := make([]string, 0, len(texts))
	for _, sentence := range texts {
		cleaned := p.cleanString(sentence)

		if cleaned == "" || strings.Contains(cleaned, "cookies") {
			continue
		}

		chars := make([]string, 0, len(cleaned))
		letterCount := 0
		for _, r := range cleaned {
			if unicode.IsLetter(r) {
				letterCount++
			}

			if (unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || unicode.IsPunct(r)) && string(r) != "-" && string(r) != ":" {
				chars = append(chars, string(r))
			}
		}
		outChars := p.cleanString(strings.Join(chars, ""))
		if letterCount < 5 || strings.Contains(outChars, "security service to protect itself from online attacks") || strings.Contains(outChars, "attention required!") {
			continue
		}

		outTexts = append(outTexts, outChars)
	}

	out := strings.Join(outTexts, " ")
	if len(out) < 100 {
		return "", NewNoDataError()
	}

	return out, nil
}

func (p parser) cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ToLower(s)
	return s
}
