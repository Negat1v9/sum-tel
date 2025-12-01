package tgparser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CleanMessageText removes HTML tags from the given text and returns plain text
func CleanMessageText(text string) string {
	// Remove html tags
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return text
	}
	return doc.Text()
}
