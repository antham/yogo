package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const refURL = "https://yopmail.com"
const defaultHttpTimeout = 10

type mailKind string

const (
	mailHTML   mailKind = "m"
	mailText            = "t"
	mailSource          = "s"
)

type MailDoc interface {
	MailHTMLDoc | MailSourceDoc | MailTextDoc
	Find(selector string) *goquery.Selection
}

type MailHTMLDoc goquery.Document
type MailSourceDoc goquery.Document
type MailTextDoc goquery.Document

// Client provides a high level interface to abstract yopmail data fetching
type Client[M MailDoc] struct {
	browser    *browser
	apiVersion string
}

// New creates a new client
func New[M MailDoc]() (Client[M], error) {
	browser := newBrowser()
	c, err := browser.get("GET", refURL, map[string]string{}, nil)
	if err != nil {
		return Client[M]{}, err
	}

	_, err = browser.get("GET", refURL+"/consent?c=accept", map[string]string{}, nil)
	if err != nil {
		return Client[M]{}, err
	}
	apiVersion, err := parseApiVersion(c.String())
	if err != nil {
		return Client[M]{}, err
	}
	return Client[M]{apiVersion: apiVersion, browser: browser}, nil
}

// GetMailsPage fetches all html pages containing emails data
func (c Client[M]) GetMailsPage(identifier string, page int) (*goquery.Document, error) {
	URL, err := decorateURL("inbox?d=&ctrl=&scrl=&spam=true&r_c=&id=", c.apiVersion, false, map[string]string{"login": identifier, "p": strconv.Itoa(page)})
	if err != nil {
		return nil, err
	}
	c.browser.populateCookieFromAccount(identifier)
	r, err := c.browser.get("GET", URL, map[string]string{}, nil)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(r)
}

// GetMailPage fetches html page containing the email
func (c Client[M]) GetMailPage(identifier string, mailID string) (doc M, err error) {
	var kind mailKind
	switch any(doc).(type) {
	case MailHTMLDoc:
		kind = mailHTML
	}

	URL, err := decorateURL("mail", c.apiVersion, true, map[string]string{"b": identifier, "id": fmt.Sprintf("%s%s", kind, mailID)})
	if err != nil {
		return
	}
	c.browser.populateCookieFromAccount(identifier)
	r, err := c.browser.get("GET", URL, map[string]string{}, nil)
	if err != nil {
		return
	}
	d, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return
	}
	return M(*d), nil
}

// DeleteMail removes an email from yopmail inbox
func (c Client[M]) DeleteMail(identifier string, mailID string) error {
	URL, err := decorateURL("inbox?p=1&ctrl=&r_c=&id=", c.apiVersion, false, map[string]string{"login": identifier, "d": mailID})
	if err != nil {
		return err
	}
	c.browser.populateCookieFromAccount(identifier)
	_, err = c.browser.get("GET", URL, map[string]string{}, nil)
	return err
}

// FlushMail removes all yopmail inbox mails
func (c Client[M]) FlushMail(identifier string, mailID string) error {
	URL, err := decorateURL("inbox?p=1&d=all&r_c=&id=", c.apiVersion, false, map[string]string{"login": identifier, "ctrl": mailID})
	if err != nil {
		return err
	}
	c.browser.populateCookieFromAccount(identifier)
	_, err = c.browser.get("GET", URL, map[string]string{}, nil)
	return err
}

func decorateURL(URL string, apiVersion string, disableDefaultQueryParams bool, queryParams map[string]string) (string, error) {
	doc, err := fetchDocument("GET", refURL, map[string]string{}, nil)
	if err != nil {
		return "", err
	}

	var yp string
	var ok bool
	doc.Find("#yp").Each(func(i int, s *goquery.Selection) {
		yp, ok = s.Attr("value")
	})
	if !ok || yp == "" {
		return "", errors.New("failure when fetching yp value")
	}

	doc, err = fetchDocument("GET", refURL+"/ver/"+apiVersion+"/webmail.js", map[string]string{}, nil)
	if err != nil {
		return "", err
	}

	m := regexp.MustCompile("&yj=(.*?)&").FindStringSubmatch(doc.Text())
	if len(m) != 2 {
		return "", errors.New("failure when fetching yj value")
	}

	yj := m[1]

	u, err := url.Parse(refURL + "/" + URL)
	if err != nil {
		return "", err
	}

	q := u.Query()

	if !disableDefaultQueryParams {
		q.Add("yp", yp)
		q.Add("yj", yj)
		q.Add("v", apiVersion)
	}
	for k, v := range queryParams {
		q.Add(k, v)
	}

	return (&url.URL{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Path:     "en" + u.Path,
		RawQuery: q.Encode(),
	}).String(), nil
}

func fetch(method string, URL string, headers map[string]string, body io.Reader) (io.Reader, error) {
	errMsg := fmt.Sprintf("failure when fetching %s", URL)
	r, err := http.NewRequest(method, URL, body)
	if err != nil {
		return nil, wrapError(errMsg, err)
	}

	for k, v := range headers {
		r.Header.Add(k, v)
	}

	c := http.Client{Timeout: defaultHttpTimeout * time.Second}
	res, err := c.Do(r)
	if err != nil {
		return nil, wrapError(errMsg, err)
	}
	if res.StatusCode > 300 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, wrapError(errMsg, errors.New("can't extract request body"))
		}

		return nil, wrapError(errMsg, fmt.Errorf("request failed with error code %d and body %s", res.StatusCode, string(b)))
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, wrapError(errMsg, err)
	}

	return bytes.NewBuffer(b), nil
}

func fetchDocument(method string, URL string, headers map[string]string, body io.Reader) (*goquery.Document, error) {
	r, err := fetch("GET", URL, headers, body)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(r)
}

func parseApiVersion(s string) (string, error) {
	data := regexp.MustCompile(`<script src="/ver/([0-9.]+)/webmail.js">`).FindStringSubmatch(s)
	if len(data) < 2 {
		return "", errors.New("api version could not be extracted")
	}

	return data[1], nil
}

func wrapError(msg string, err error) error {
	return fmt.Errorf("%s : %w", msg, err)
}

type browser struct {
	cookies map[string]string
}

func newBrowser() *browser {
	return &browser{cookies: map[string]string{}}
}

func (b *browser) setCookie(key string, value string) {
	b.cookies[key] = value
}

func (b *browser) populateCookieFromAccount(account string) {
	for k, v := range map[string]string{"compte": account, "ywm": account, "ytime": time.Now().Format("15:04")} {
		b.setCookie(k, v)
	}
}

func (b *browser) buildCookie() string {
	data := []string{}
	for k, v := range b.cookies {
		data = append(data, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(data, "; ")
}

func (b *browser) get(method string, URL string, headers map[string]string, body io.Reader) (*bytes.Buffer, error) {
	errMsg := fmt.Sprintf("failure when fetching %s", URL)
	r, err := http.NewRequest(method, URL, body)
	if err != nil {
		return nil, wrapError(errMsg, err)
	}

	for k, v := range headers {
		r.Header.Add(k, v)
	}
	if len(b.cookies) > 0 {
		r.Header.Add("Cookie", b.buildCookie())
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")

	c := http.Client{Timeout: defaultHttpTimeout * time.Second}
	res, err := c.Do(r)
	if err != nil {
		return nil, wrapError(errMsg, err)
	}
	if res.StatusCode > 300 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, wrapError(errMsg, errors.New("can't extract request body"))
		}

		return nil, wrapError(errMsg, fmt.Errorf("request failed with error code %d and body %s", res.StatusCode, string(b)))
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, wrapError(errMsg, err)
	}

	for _, c := range res.Cookies() {
		b.cookies[c.Name] = c.Value
	}

	return bytes.NewBuffer(buf), nil
}
