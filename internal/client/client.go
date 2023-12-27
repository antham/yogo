package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"golang.org/x/net/http/httpproxy"
)

var ErrCaptcha = errors.New("failure when trying to access content: a CAPTCHA is probably activated, look to the web interface")

const refURL = "https://yopmail.com"
const defaultRequestTimeout = 10
const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"

type mailKind string

const (
	mailHTML   mailKind = "m"
	mailText   mailKind = "t"
	mailSource mailKind = "s"
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
func New[M MailDoc](enableDebugMode bool) (Client[M], error) {
	browser := newBrowser(enableDebugMode)
	c, err := browser.fetch("GET", refURL, map[string]string{}, nil)
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
	URL, err := c.decorateURL("inbox?d=&ctrl=&scrl=&spam=true&ad=0&r_c=&id=", c.apiVersion, false, map[string]string{"login": identifier, "p": strconv.Itoa(page)})
	if err != nil {
		return nil, err
	}
	c.browser.populateCookieFromAccount(identifier)
	content, err := c.browser.fetch("GET", URL, map[string]string{}, nil)
	if err != nil {
		return nil, err
	}
	err = checkInboxCAPTCHA(content.String())
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(content)
}

// GetMailPage fetches html page containing the email
func (c Client[M]) GetMailPage(identifier string, mailID string) (doc M, err error) {
	var kind mailKind
	switch any(doc).(type) {
	case MailHTMLDoc:
		kind = mailHTML
	case MailSourceDoc:
		kind = mailSource
	}
	URL, err := c.decorateURL("mail", c.apiVersion, true, map[string]string{"b": identifier, "id": fmt.Sprintf("%s%s", kind, mailID)})
	if err != nil {
		return
	}
	c.browser.populateCookieFromAccount(identifier)
	content, err := c.browser.fetch("GET", URL, map[string]string{}, nil)
	if err != nil {
		return
	}
	err = checkMailCAPTCHA(content.String())
	if err != nil {
		return
	}
	d, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return
	}
	return M(*d), nil
}

// DeleteMail removes an email from yopmail inbox
func (c Client[M]) DeleteMail(identifier string, mailID string) error {
	URL, err := c.decorateURL("inbox?p=1&ctrl=&ad=0&r_c=&id=", c.apiVersion, false, map[string]string{"login": identifier, "d": mailID})
	if err != nil {
		return err
	}
	c.browser.populateCookieFromAccount(identifier)
	content, err := c.browser.fetch("GET", URL, map[string]string{}, nil)
	if err != nil {
		return err
	}
	return checkInboxCAPTCHA(content.String())
}

// FlushMail removes all yopmail inbox mails
func (c Client[M]) FlushMail(identifier string, mailID string) error {
	URL, err := c.decorateURL("inbox?p=1&d=all&ad=0&r_c=&id=", c.apiVersion, false, map[string]string{"login": identifier, "ctrl": mailID})
	if err != nil {
		return err
	}
	c.browser.populateCookieFromAccount(identifier)
	content, err := c.browser.fetch("GET", URL, map[string]string{}, nil)
	if err != nil {
		return err
	}
	return checkInboxCAPTCHA(content.String())
}

func (c Client[M]) decorateURL(URL string, apiVersion string, disableDefaultQueryParams bool, queryParams map[string]string) (string, error) {
	doc, err := c.browser.fetchDocument("GET", refURL, map[string]string{}, nil)
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
	doc, err = c.browser.fetchDocument("GET", refURL+"/ver/"+apiVersion+"/webmail.js", map[string]string{}, nil)
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

type browser struct {
	cookies           map[string]string
	enableDebugMode   bool
	httpClientFactory httpClientFactory
}

func newBrowser(enableDebugMode bool) *browser {
	return &browser{
		cookies:           map[string]string{},
		enableDebugMode:   enableDebugMode,
		httpClientFactory: httpClientFactory{},
	}
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

func (b *browser) fetch(method string, URL string, headers map[string]string, body io.Reader) (*bytes.Buffer, error) {
	ID := uuid.New()
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
	userAgent := os.Getenv("YOGO_USER_AGENT")
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	r.Header.Add("User-Agent", userAgent)
	if b.enableDebugMode {
		buff, err := httputil.DumpRequest(r, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("\n---- REQUEST %s ----\n", ID)
		fmt.Println(string(buff))
		fmt.Println("------------------------------------------------------")
	}

	c, err := b.httpClientFactory.create()
	if err != nil {
		return nil, wrapError(errMsg, err)
	}
	res, err := c.Do(r)
	if err != nil {
		return nil, wrapError(errMsg, err)
	}
	if b.enableDebugMode {
		buff, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
		fmt.Printf("\n---- RESPONSE %s ----\n", ID)
		fmt.Println(string(buff))
		fmt.Println("-------------------------------------------------------")
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

func (b *browser) fetchDocument(method string, URL string, headers map[string]string, body io.Reader) (*goquery.Document, error) {
	r, err := b.fetch("GET", URL, headers, body)
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

func checkInboxCAPTCHA(content string) error {
	s := `w\.finrmail\(\d+,\s*\d+,\s*\d+,\s*\d+,\s*\d+,\s*'alt\.[^']+',\s*'.*?'\)|Loading \.\.\.`
	if !regexp.MustCompile(s).MatchString(content) {
		return ErrCaptcha
	}
	return nil
}

func checkMailCAPTCHA(content string) error {
	if strings.Contains(content, "window.showRc()") {
		return ErrCaptcha
	}
	return nil
}

type httpClientFactory struct{}

func (h httpClientFactory) create() (*http.Client, error) {
	timeout := defaultRequestTimeout * time.Second
	if os.Getenv("YOGO_REQUEST_TIMEOUT") != "" {
		t, err := strconv.Atoi(os.Getenv("YOGO_REQUEST_TIMEOUT"))
		if err != nil {
			return nil, err
		}
		timeout = time.Duration(t) * time.Second
	}
	client := &http.Client{Timeout: timeout}
	config := httpproxy.FromEnvironment()
	URL := ""
	switch {
	case config.HTTPProxy != "":
		URL = config.HTTPProxy
	case config.HTTPSProxy != "":
		URL = config.HTTPSProxy
	}
	if URL != "" {
		u, err := url.Parse(URL)
		if err != nil {
			return nil, err
		}
		transport := &http.Transport{}
		transport.Proxy = http.ProxyURL(u)
		client.Transport = transport
	}
	return client, nil
}
