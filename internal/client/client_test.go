package client

import (
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestParseApiVersion(t *testing.T) {
	type scenario struct {
		name      string
		arguments func() string
		test      func(string, error)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, s := range []scenario{{
		"no version found in JS file",
		func() string {
			return ""
		}, func(version string, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, "api version could not be extracted")
		},
	}, {
		"version found in JS file",
		func() string {
			return `xxxxxx<script src="/ver/3.1/webmail.js">xxxxx`
		}, func(version string, err error) {
			assert.NoError(t, err)
			assert.Equal(t, "3.1", version)
		},
	}} {
		t.Run(s.name, func(t *testing.T) {
			s.test(parseApiVersion(s.arguments()))
			httpmock.Reset()
		})
	}
}

func TestFetchDocument(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		test  func(*goquery.Document, error)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, s := range []scenario{
		{
			"error when fetching the request",
			func() {
				httpmock.RegisterResponder("GET", "http://abcdefg.com",
					httpmock.NewErrorResponder(errors.New("an error occurred")))
			}, func(content *goquery.Document, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, `failure when fetching http://abcdefg.com : Get "http://abcdefg.com": an error occurred`)
			},
		}, {
			"fetch an html document",
			func() {
				httpmock.RegisterResponder("GET", "http://abcdefg.com",
					httpmock.NewStringResponder(200, `<html><head></head><body><div></div></body></html>`))
			}, func(content *goquery.Document, err error) {
				assert.NoError(t, err)
				doc, err := content.Html()
				assert.NoError(t, err)
				assert.Equal(t, "<html><head></head><body><div></div></body></html>", doc)
			},
		}} {
		t.Run(s.name, func(t *testing.T) {
			b := newBrowser()

			s.setup()
			s.test(b.fetchDocument("GET", "http://abcdefg.com", map[string]string{}, nil))
			httpmock.Reset()
		})
	}
}

func TestFetch(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		test  func(io.Reader, error)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, s := range []scenario{{
		"500 on response",
		func() {
			httpmock.RegisterResponder("GET", "http://hijklm.com",
				httpmock.NewStringResponder(500, ""))
		}, func(reader io.Reader, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when fetching http://hijklm.com : request failed with error code 500 and body `)
		},
	}, {
		"error when fetching the request",
		func() {
			httpmock.RegisterResponder("GET", "http://hijklm.com",
				httpmock.NewErrorResponder(errors.New("an error occurred")))
		}, func(reader io.Reader, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when fetching http://hijklm.com : Get "http://hijklm.com": an error occurred`)
		},
	}, {
		"request timeout",
		func() {
			httpmock.RegisterResponder("GET", "http://hijklm.com",
				func(r *http.Request) (res *http.Response, e error) {
					time.Sleep(time.Second * 20)
					return
				})
		}, func(reader io.Reader, err error) {
			assert.Error(t, err)
			assert.Regexp(t, "lient.Timeout exceeded while awaiting headers|context deadline exceeded", err.Error())
		},
	}, {
		"fetch an URL",
		func() {
			httpmock.RegisterResponder("GET", "http://hijklm.com",
				httpmock.NewStringResponder(200, `<html><head></head><body><div></div></body></html>`))
		}, func(reader io.Reader, err error) {
			assert.NoError(t, err)
			b, err := io.ReadAll(reader)
			assert.NoError(t, err)
			assert.Equal(t, "<html><head></head><body><div></div></body></html>", string(b))
		},
	}} {
		t.Run(s.name, func(t *testing.T) {
			b := newBrowser()
			s.setup()
			s.test(b.fetch("GET", "http://hijklm.com", map[string]string{"header1": "value1", "header2": "value2"}, nil))
			httpmock.Reset()
		})
	}
}

func TestDecorateURL(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		args  func() (string, string, bool, map[string]string)
		test  func(string, error)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, s := range []scenario{{
		"500 when requesting yopmail",
		func() {
			httpmock.RegisterResponder("GET", refURL,
				httpmock.NewStringResponder(500, ""))
		}, func() (string, string, bool, map[string]string) {
			return "test", "3.1", false, map[string]string{"q1": "value1", "q2": "value2"}
		},
		func(URL string, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when fetching https://yopmail.com : request failed with error code 500 and body `)
		},
	}, {
		"no attribute yp found",
		func() {
			httpmock.RegisterResponder("GET", refURL,
				httpmock.NewStringResponder(200, ""))
		}, func() (string, string, bool, map[string]string) {
			return "test", "3.1", false, map[string]string{"q1": "value1", "q2": "value2"}
		}, func(URL string, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, "failure when fetching yp value")
		},
	}, {
		"attribute yp with no value",
		func() {
			httpmock.RegisterResponder("GET", refURL,
				httpmock.NewStringResponder(200, `<html><head></head><body><input id="yp"></body></html>`))
		}, func() (string, string, bool, map[string]string) {
			return "test", "3.1", false, map[string]string{"q1": "value1", "q2": "value2"}
		}, func(URL string, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, "failure when fetching yp value")
		},
	}, {
		"failure when fetching the JS file",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/ver/3.1/webmail.js",
				httpmock.NewStringResponder(500, ""))
			httpmock.RegisterResponder("GET", refURL,
				httpmock.NewStringResponder(200, `<html><head></head><body><input id="yp" value="yptest"></body></html>`))
		}, func() (string, string, bool, map[string]string) {
			return "test", "3.1", false, map[string]string{"q1": "value1", "q2": "value2"}
		}, func(URL string, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, "failure when fetching https://yopmail.com/ver/3.1/webmail.js : request failed with error code 500 and body ")
		},
	}, {
		"no yj attribute",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/ver/3.1/webmail.js",
				httpmock.NewStringResponder(200, ""))
			httpmock.RegisterResponder("GET", refURL,
				httpmock.NewStringResponder(200, `<html><head></head><body><input id="yp" value="yptest"></body></html>`))
		}, func() (string, string, bool, map[string]string) {
			return "test", "3.1", false, map[string]string{"q1": "value1", "q2": "value2"}
		}, func(URL string, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, "failure when fetching yj value")
		},
	}, {
		"failure when parsing the URL",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/ver/3.1/webmail.js",
				httpmock.NewStringResponder(200, "xxx http://whatever.com?q=s&yj=ytest&t=a xxxxx"))
			httpmock.RegisterResponder("GET", refURL,
				httpmock.NewStringResponder(200, `<html><head></head><body><input id="yp" value="yptest"</body></html>`))
		}, func() (string, string, bool, map[string]string) {
			return "\n\n", "3.1", false, map[string]string{"q1": "value1", "q2": "value2"}
		}, func(URL string, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `parse "https://yopmail.com/\n\n": net/url: invalid control character in URL`)
		},
	}, {
		"decorate the URL",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/ver/3.1/webmail.js",
				httpmock.NewStringResponder(200, "xxx http://whatever.com?q=s&yj=ytest&t=a xxxxx"))
			httpmock.RegisterResponder("GET", refURL,
				httpmock.NewStringResponder(200, `<html><head></head><body><input id="yp" value="yptest"></body></html>`))
		}, func() (string, string, bool, map[string]string) {
			return "test?k=w&g=t", "3.1", false, map[string]string{"q1": "value1", "q2": "value2"}
		}, func(URL string, err error) {
			assert.NoError(t, err)
			assert.Equal(t, refURL+"/en/test?g=t&k=w&q1=value1&q2=value2&v=3.1&yj=ytest&yp=yptest", URL)
		},
	}, {
		"decorate the URL and do not add default query params",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/ver/3.1/webmail.js",
				httpmock.NewStringResponder(200, "xxx http://whatever.com?q=s&yj=ytest&t=a xxxxx"))
			httpmock.RegisterResponder("GET", refURL,
				httpmock.NewStringResponder(200, `<html><head></head><body><input id="yp" value="yptest"></body></html>`))
		}, func() (string, string, bool, map[string]string) {
			return "test?k=w&g=t", "3.1", true, map[string]string{"q1": "value1", "q2": "value2"}
		}, func(URL string, err error) {
			assert.NoError(t, err)
			assert.Equal(t, refURL+"/en/test?g=t&k=w&q1=value1&q2=value2", URL)
		},
	}} {
		t.Run(s.name, func(t *testing.T) {
			mockYopmailSetup()

			c, err := New[MailHTMLDoc]()
			assert.NoError(t, err)

			s.setup()
			s.test(c.decorateURL(s.args()))
			httpmock.Reset()
		})
	}
}

func TestGetMailsPage(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		args  func() (string, int)
		test  func(*goquery.Document, error)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, s := range []scenario{{
		"500 when requesting yopmail",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/inbox?ctrl=&d=&id=&login=box1&p=1&r_c=&scrl=&spam=true&v=3.1&yj=ytest&yp=yptest",
				httpmock.NewStringResponder(500, ""))
		}, func() (string, int) {
			return "box1", 1
		},
		func(doc *goquery.Document, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when fetching https://yopmail.com/en/inbox?ctrl=&d=&id=&login=box1&p=1&r_c=&scrl=&spam=true&v=3.1&yj=ytest&yp=yptest : request failed with error code 500 and body `)
		},
	}, {
		"CAPTCHA activated",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/inbox?ctrl=&d=&id=&login=box1&p=1&r_c=&scrl=&spam=true&v=3.1&yj=ytest&yp=yptest",
				httpmock.NewStringResponder(200, ""))
		}, func() (string, int) {
			return "box1", 1
		},
		func(doc *goquery.Document, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when trying to access content: a CAPTCHA is probably activated, look to the web interface`)
		},
	}, {
		"request succeed",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/inbox?ctrl=&d=&id=&login=box1&p=1&r_c=&scrl=&spam=true&v=3.1&yj=ytest&yp=yptest",
				httpmock.NewStringResponder(200, "w.finrmail(25,2,1,0,0,'alt.zk-4nyqp5l','')"))
		}, func() (string, int) {
			return "box1", 1
		},
		func(doc *goquery.Document, err error) {
			assert.NoError(t, err)
		},
	}} {
		t.Run(s.name, func(t *testing.T) {
			mockYopmailSetup()

			c, err := New[MailHTMLDoc]()
			assert.NoError(t, err)

			s.setup()
			s.test(c.GetMailsPage(s.args()))
			httpmock.Reset()
		})
	}
}

func TestGetMailPage(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		args  func() (string, string)
		test  func(MailHTMLDoc, error)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, s := range []scenario{{
		"500 when requesting yopmail",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/mail?b=box1&id=mABCDEFGH",
				httpmock.NewStringResponder(500, ""))
		}, func() (string, string) {
			return "box1", "ABCDEFGH"
		},
		func(doc MailHTMLDoc, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when fetching https://yopmail.com/en/mail?b=box1&id=mABCDEFGH : request failed with error code 500 and body `)
		},
	}, {
		"CAPTCHA activated",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/mail?b=box1&id=mABCDEFGH",
				httpmock.NewStringResponder(200, "window.showRc()"))
		}, func() (string, string) {
			return "box1", "ABCDEFGH"
		},
		func(doc MailHTMLDoc, err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when trying to access content: a CAPTCHA is probably activated, look to the web interface`)
		},
	}, {
		"request succeed",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/mail?b=box1&id=mABCDEFGH",
				httpmock.NewStringResponder(200, "w.finrmail(25,2,1,0,0,'alt.zk-4nyqp5l','')"))
		}, func() (string, string) {
			return "box1", "ABCDEFGH"
		},
		func(doc MailHTMLDoc, err error) {
			assert.NoError(t, err)
		},
	}} {
		t.Run(s.name, func(t *testing.T) {
			mockYopmailSetup()

			c, err := New[MailHTMLDoc]()
			assert.NoError(t, err)

			s.setup()
			s.test(c.GetMailPage(s.args()))
			httpmock.Reset()
		})
	}
}

func TestDeleteMail(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		args  func() (string, string)
		test  func(error)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, s := range []scenario{{
		"500 when requesting yopmail",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/inbox?ctrl=&d=ABCDEFGH&login=box1&p=1&r_c=&id=&v=3.1&yj=ytest&yp=yptest",
				httpmock.NewStringResponder(500, ""))
		}, func() (string, string) {
			return "box1", "ABCDEFGH"
		},
		func(err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when fetching https://yopmail.com/en/inbox?ctrl=&d=ABCDEFGH&id=&login=box1&p=1&r_c=&v=3.1&yj=ytest&yp=yptest : Get "https://yopmail.com/en/inbox?ctrl=&d=ABCDEFGH&id=&login=box1&p=1&r_c=&v=3.1&yj=ytest&yp=yptest": no responder found`)
		},
	}, {
		"CAPTCHA activated",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/inbox?ctrl=&d=ABCDEFGH&id=&login=box1&p=1&r_c=&v=3.1&yj=ytest&yp=yptest",
				httpmock.NewStringResponder(200, ""))
		}, func() (string, string) {
			return "box1", "ABCDEFGH"
		},
		func(err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when trying to access content: a CAPTCHA is probably activated, look to the web interface`)
		},
	}, {
		"request succeed",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/inbox?ctrl=&d=ABCDEFGH&id=&login=box1&p=1&r_c=&v=3.1&yj=ytest&yp=yptest",
				httpmock.NewStringResponder(200, "w.finrmail(25,2,1,0,0,'alt.zk-4nyqp5l','')"))
		}, func() (string, string) {
			return "box1", "ABCDEFGH"
		},
		func(err error) {
			assert.NoError(t, err)
		},
	}} {
		t.Run(s.name, func(t *testing.T) {
			mockYopmailSetup()

			c, err := New[MailHTMLDoc]()
			assert.NoError(t, err)

			s.setup()
			s.test(c.DeleteMail(s.args()))
			httpmock.Reset()
		})
	}
}

func TestFlushMail(t *testing.T) {
	type scenario struct {
		name  string
		setup func()
		args  func() (string, string)
		test  func(error)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, s := range []scenario{{
		"500 when requesting yopmail",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/inbox?ctrl=ABCDEFGH&d=all&id=&login=box1&p=1&r_c=&v=3.1&yj=ytest&yp=yptest",
				httpmock.NewStringResponder(500, ""))
		}, func() (string, string) {
			return "box1", "ABCDEFGH"
		},
		func(err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when fetching https://yopmail.com/en/inbox?ctrl=ABCDEFGH&d=all&id=&login=box1&p=1&r_c=&v=3.1&yj=ytest&yp=yptest : request failed with error code 500 and body `)
		},
	}, {
		"CAPTCHA activated",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/inbox?ctrl=ABCDEFGH&d=all&id=&login=box1&p=1&r_c=&v=3.1&yj=ytest&yp=yptest",
				httpmock.NewStringResponder(200, ""))
		}, func() (string, string) {
			return "box1", "ABCDEFGH"
		},
		func(err error) {
			assert.Error(t, err)
			assert.EqualError(t, err, `failure when trying to access content: a CAPTCHA is probably activated, look to the web interface`)
		},
	}, {
		"request succeed",
		func() {
			httpmock.RegisterResponder("GET", refURL+"/en/inbox?ctrl=ABCDEFGH&d=all&id=&login=box1&p=1&r_c=&v=3.1&yj=ytest&yp=yptest",
				httpmock.NewStringResponder(200, "w.finrmail(25,2,1,0,0,'alt.zk-4nyqp5l','')"))
		}, func() (string, string) {
			return "box1", "ABCDEFGH"
		},
		func(err error) {
			assert.NoError(t, err)
		},
	}} {
		t.Run(s.name, func(t *testing.T) {
			mockYopmailSetup()

			c, err := New[MailHTMLDoc]()
			assert.NoError(t, err)

			s.setup()
			s.test(c.FlushMail(s.args()))
			httpmock.Reset()
		})
	}
}

func mockYopmailSetup() {
	httpmock.RegisterResponder("GET", refURL+"/ver/3.1/webmail.js",
		httpmock.NewStringResponder(200, "xxx http://whatever.com?q=s&yj=ytest&t=a xxxxx"))
	httpmock.RegisterResponder("GET", refURL+"/consent?c=accept", httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("GET", refURL,
		httpmock.NewStringResponder(200, `<script src="/ver/3.1/webmail.js"></script><input id="yp" value="yptest">`))
}
