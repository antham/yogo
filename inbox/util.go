package inbox

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
)

var send = func(URL string) error {
	_, err := http.Get(URL)
	return err
}

var buildReader = func(method string, URL string, headers map[string]string, body io.Reader) (io.Reader, error) {
	r, err := http.NewRequest(method, URL, body)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		r.Header.Add(k, v)
	}

	c := http.Client{}
	res, err := c.Do(r)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(b), nil
}

var fetchFromReader = func(r io.Reader) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return nil, err
	}

	return doc, err
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
