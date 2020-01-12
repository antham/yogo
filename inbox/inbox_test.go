package inbox

import (
	"io/ioutil"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestParseMailID(t *testing.T) {
	ID := parseMailID(`http://www.yopmail.com/m.php?b=test&id=me_ZGLjAmN5ZGLmZwVlZQNjZmH3BQN5ZD==`)

	assert.Equal(t, ID, "me_ZGLjAmN5ZGLmZwVlZQNjZmH3BQN5ZD==")
}

func TestParseInboxPage(t *testing.T) {
	inbox := &Inbox{}

	parseInboxPage(getDoc("inbox_page.html"), inbox)

	assert.Equal(t, 15, inbox.Count())
	assert.Equal(t, "me_ZGtkZwVmZQNlZmVlZQNjZQRlAwHlZD==", inbox.Get(0).ID)
	assert.Equal(t, "me_ZGtkZwVmZQNkBGDkZQNjZQRkZQZ0AN==", inbox.Get(14).ID)
}

func TestCount(t *testing.T) {
	inbox := &Inbox{}

	parseInboxPage(getDoc("inbox_page.html"), inbox)

	assert.Equal(t, inbox.Count(), 15)
}

func TestParseInboxPages(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_1.html",
		},
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_2.html",
		},
		{
			"GET",
			"http://www.yopmail.com",
			"main_page.html",
		},
		{
			"GET",
			"http://www.yopmail.com/style/2.9/webmail.js",
			"webmail.js",
		},
	}))

	inbox, err := ParseInboxPages("test", 29)

	assert.NoError(t, err)
	assert.Equal(t, "test", inbox.GetIdentifier())
	assert.Equal(t, 29, inbox.Count())
	assert.Equal(t, "me_ZGtkZwVmZQNmBGV1ZQNjZQVjAwD1BD==", inbox.Get(0).ID)
	assert.Equal(t, "me_ZGtkZwVmZQNmAQH3ZQNjZQR4AGVlAD==", inbox.Get(28).ID)
}

func TestShrink(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_1.html",
		},
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_2.html",
		},
		{
			"GET",
			"http://www.yopmail.com",
			"main_page.html",
		},
		{
			"GET",
			"http://www.yopmail.com/style/2.9/webmail.js",
			"webmail.js",
		},
	}))

	inbox, err := ParseInboxPages("test", 19)

	assert.NoError(t, err)
	assert.Equal(t, 19, inbox.Count())
	assert.Equal(t, "me_ZGtkZwVmZQNmBGV1ZQNjZQVjAwD1BD==", inbox.Get(0).ID)
	assert.Equal(t, "me_ZGtkZwVmZQNmAwRjZQNjZQR5ZQp4BN==", inbox.Get(18).ID)
}

func TestShrinkEmptyInbox(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_empty.html",
		},
		{
			"GET",
			"http://www.yopmail.com",
			"main_page.html",
		},
		{
			"GET",
			"http://www.yopmail.com/style/2.9/webmail.js",
			"webmail.js",
		},
	}))

	inbox, err := ParseInboxPages("test", 1)

	assert.NoError(t, err)
	assert.Equal(t, 0, inbox.Count())
}

func TestShrinkWithLimitGreaterThanNumberOfMessagesAvailable(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_1.html",
		},
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_empty.html",
		},
		{
			"GET",
			"http://www.yopmail.com",
			"main_page.html",
		},
		{
			"GET",
			"http://www.yopmail.com/style/2.9/webmail.js",
			"webmail.js",
		},
	}))

	inbox, err := ParseInboxPages("test", 18)

	assert.NoError(t, err)
	assert.Equal(t, 15, inbox.Count())
}

func TestGetAll(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_1.html",
		},
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_2.html",
		},
		{
			"GET",
			"http://www.yopmail.com",
			"main_page.html",
		},
		{
			"GET",
			"http://www.yopmail.com/style/2.9/webmail.js",
			"webmail.js",
		},
	}))

	inbox, err := ParseInboxPages("test", 29)
	mails := inbox.GetAll()

	assert.NoError(t, err)
	assert.Len(t, mails, 29)
	assert.Equal(t, "me_ZGtkZwVmZQNmBGV1ZQNjZQVjAwD1BD==", mails[0].ID)
	assert.Equal(t, "me_ZGtkZwVmZQNmAQH3ZQNjZQR4AGVlAD==", mails[28].ID)
}

func TestFlush(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_1.html",
		},
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_2.html",
		},
		{
			"GET",
			"http://www.yopmail.com",
			"main_page.html",
		},
		{
			"GET",
			"http://www.yopmail.com/style/2.9/webmail.js",
			"webmail.js",
		},
		{
			"GET",
			"http://www.yopmail.com/inbox.php?login=test&p=1&d=all&ctrl=e_ZGtkZwVmZQNmBGV1ZQNjZQVjAwD1BD==&v=2.9&r_c=&id",
			"noop.html",
		},
	}))

	inbox, err := ParseInboxPages("test", 15)
	inbox.Flush()

	assert.NoError(t, err)
}

func TestFlushEmptyInbox(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_empty.html",
		},
		{
			"GET",
			"http://www.yopmail.com",
			"main_page.html",
		},
		{
			"GET",
			"http://www.yopmail.com/style/2.9/webmail.js",
			"webmail.js",
		},
	}))

	inbox, err := ParseInboxPages("test", 1)
	inbox.Flush()

	assert.NoError(t, err)
	assert.Equal(t, 0, inbox.Count())
}

func TestDelete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"http://www.yopmail.com/inbox.php?ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=2.9&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UZGx3Zmt3AQL1ZmZ1ZQV1Zwx",
			"inbox_page_1.html",
		},
		{
			"GET",
			"http://www.yopmail.com",
			"main_page.html",
		},
		{
			"GET",
			"http://www.yopmail.com/style/2.9/webmail.js",
			"webmail.js",
		},
		{
			"GET",
			"http://www.yopmail.com/inbox.php?login=test&p=1&d=e_ZGtkZwVmZQNmBGV1ZQNjZQVjAwD1BD==&ctrl=&scrl=0&spam=true&v=2.9&r_c=",
			"noop.html",
		},
	}))

	inbox, err := ParseInboxPages("test", 1)
	assert.NoError(t, inbox.Delete(0))
	assert.Equal(t, 1, httpmock.GetCallCountInfo()["GET http://www.yopmail.com/inbox.php?login=test&p=1&d=e_ZGtkZwVmZQNmBGV1ZQNjZQVjAwD1BD==&ctrl=&scrl=0&spam=true&v=2.9&r_c="])
	assert.NoError(t, err)
}

type responder struct {
	method   string
	URL      string
	filename string
}

func registerResponders(responders []responder) error {
	for _, r := range responders {
		b, err := ioutil.ReadFile(r.filename)
		if err != nil {
			return err
		}

		httpmock.RegisterResponder(r.method, r.URL,
			httpmock.NewStringResponder(200, string(b)))
	}
	return nil
}
