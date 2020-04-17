package inbox

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
)

var apiVersion = ""

const refURL = "http://www.yopmail.com"

func init() {
	fetchApiVersion()
}

func send(URL string) error {
	_, err := http.Get(URL)
	return err
}

func buildReader(method string, URL string, headers map[string]string, body io.Reader) (io.Reader, error) {
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

func fetchFromReader(r io.Reader) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(r)
}

func fetchURL(URL string) (*goquery.Document, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromReader(resp.Body)
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

func fetchApiVersion() {
	if apiVersion != "" {
		return
	}

	res, err := http.Get(refURL)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := regexp.MustCompile(`<script type="text/javascript" src="/style/([0-9.]+)/webmail.js">`).FindStringSubmatch(string(b))

	if len(data) < 2 {
		log.Fatal("no yopmail api version found")
	}

	apiVersion = data[1]
}
