package inbox

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
)

var send = func(URL string) {
	http.Get(URL)
}

var fetchURL = func(URL string) (*goquery.Document, error) {
	doc, err := goquery.NewDocument(URL)

	if err != nil {
		return nil, err
	}

	return doc, err
}

func parseHTML(content string, err error) string {
	if err != nil {
		return ""
	}

	text, err := html2text.FromString(content)

	if err != nil {
		return ""
	}

	return text
}
