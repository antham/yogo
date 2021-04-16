package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const refURL = "http://www.yopmail.com"
const defaultHttpTimeout = 10

// Client provides a high level interface to abstract yopmail data fetching
type Client struct {
	apiVersion string
}

// New creates a new client
func New() (Client, error) {
	apiVersion, err := fetchApiVersion()
	if err != nil {
		return Client{}, err
	}
	return Client{apiVersion: apiVersion}, nil
}

// GetMailsPage fetches all html pages containing emails data
func (c Client) GetMailsPage(identifier string, page int) (*goquery.Document, error) {
	URL, err := decorateURL("inbox.php?d=&ctrl=&scrl=&spam=true&r_c=&id=", c.apiVersion, map[string]string{"login": identifier, "p": strconv.Itoa(page)})
	if err != nil {
		return nil, err
	}
	return fetchDocument("GET", URL, createCompteCookie(identifier), nil)
}

// GetMailPage fetches html page containing the email
func (c Client) GetMailPage(identifier string, mailID string) (*goquery.Document, error) {
	URL, err := decorateURL("m.php", c.apiVersion, map[string]string{"b": identifier, "id": mailID})
	if err != nil {
		return nil, err
	}
	return fetchDocument("GET", URL, createCompteCookie(identifier), nil)
}

// DeleteMail removes an email from yopmail inbox
func (c Client) DeleteMail(identifier string, mailID string) error {
	URL, err := decorateURL("inbox.php?p=1&ctrl=&scrl=0&spam=true&r_c=", c.apiVersion, map[string]string{"login": identifier, "d": strings.TrimLeft(mailID, "m")})
	if err != nil {
		return err
	}
	_, err = fetch("GET", URL, createCompteCookie(identifier), nil)
	return err
}

// FlushMail removes all yopmail inbox mails
func (c Client) FlushMail(identifier string, mailID string) error {
	URL, err := decorateURL("inbox.php?p=1&d=all&r_c=&id=none&scrl=&spam=true", c.apiVersion, map[string]string{"login": identifier, "ctrl": strings.TrimLeft(mailID, "m")})
	if err != nil {
		return err
	}
	_, err = fetch("GET", URL, createCompteCookie(identifier), nil)
	return err
}

func decorateURL(URL string, apiVersion string, queryParams map[string]string) (string, error) {
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

	doc, err = fetchDocument("GET", refURL+"/style/"+apiVersion+"/webmail.js", map[string]string{}, nil)
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
	q.Add("yp", yp)
	q.Add("yj", yj)
	q.Add("v", apiVersion)
	for k, v := range queryParams {
		q.Add(k, v)
	}

	return (&url.URL{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Path:     u.Path,
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
		return nil, wrapError(errMsg, fmt.Errorf("request failed with error code %d", res.StatusCode))
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

func fetchApiVersion() (string, error) {
	errMsg := "no yopmail api version found"
	client := http.Client{Timeout: defaultHttpTimeout * time.Second}
	res, err := client.Get(refURL)
	if err != nil {
		return "", wrapError(errMsg, err)
	}
	if res.StatusCode > 300 {
		return "", wrapError(errMsg, fmt.Errorf("request failed with error code %d", res.StatusCode))
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", wrapError(errMsg, err)
	}

	data := regexp.MustCompile(`<script type="text/javascript" src="/style/([0-9.]+)/webmail.js">`).FindStringSubmatch(string(b))
	if len(data) < 2 {
		return "", wrapError(errMsg, errors.New("version could not be extracted"))
	}

	return data[1], nil
}

func createCompteCookie(compte string) map[string]string {
	return map[string]string{"Cookie": fmt.Sprintf("compte=%s", compte)}
}

func wrapError(msg string, err error) error {
	return fmt.Errorf("%s : %w", msg, err)
}
