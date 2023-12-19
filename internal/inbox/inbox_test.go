package inbox

import (
	"errors"
	"os"
	"testing"

	"github.com/antham/yogo/v4/internal/client"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_1.html",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_2.html",
		},
		{
			"GET",
			"https://yopmail.com/en/mail?b=test&id=me_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj%3D%3D",
			"features/mail.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)
	err = inbox.ParseInboxPages(15)
	assert.NoError(t, err)

	m, err := inbox.Fetch(0)
	assert.NoError(t, err)
	j, err := m.JSON()
	assert.NoError(t, err)
	assert.Contains(t, j, "e_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj==")
}

func TestFetchWithDebug(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_1.html",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_2.html",
		},
		{
			"GET",
			"https://yopmail.com/en/mail?b=test&id=me_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj%3D%3D",
			"features/mail.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", true)
	assert.NoError(t, err)
	err = inbox.ParseInboxPages(15)
	assert.NoError(t, err)

	m, err := inbox.Fetch(0)
	assert.NoError(t, err)
	j, err := m.JSON()
	assert.NoError(t, err)
	assert.Contains(t, j, "e_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj==")
}

func TestCount(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_1.html",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_2.html",
		},
		{
			"GET",
			"https://yopmail.com/en/mail?b=test&id=me_ZwRjAwRkZwRkAQV1ZQNjBGD4AGL4AD%3D%3D",
			"features/mail.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)
	err = inbox.ParseInboxPages(15)
	assert.NoError(t, err)
	assert.Equal(t, inbox.Count(), 15)
}

func TestParseInboxPages(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_1.html",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_2.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/en/mail?b=test&id=me_ZwRjAwRmZGtmZQR0ZQNjAwt2BGV5BN%3D%3D",
			"features/mail.html",
		},
		{
			"GET",
			"https://yopmail.com/en/mail?b=test&id=me_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj%3D%3D",
			"features/mail.html",
		},
		{
			"GET",
			"https://yopmail.com/en/mail?b=test&id=me_ZwRjAwRmZGtmZwR0ZQNjAwt3AmxlZN%3D%3D",
			"features/mail.html",
		},
		{
			"GET",
			"https://yopmail.com/en/mail?b=test&id=me_ZwRjAwRmZGtmZwN3ZQNjAwt3AmZlAD%3D%3D",
			"features/mail2.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)

	err = inbox.ParseInboxPages(29)

	assert.NoError(t, err)
	assert.Equal(t, "test", inbox.Name)
	assert.Equal(t, 29, inbox.Count())
	m, err := inbox.Fetch(0)
	assert.NoError(t, err)
	j, err := m.JSON()
	assert.NoError(t, err)
	assert.Contains(t, j, "e_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj==")
	m, err = inbox.Fetch(28)
	assert.NoError(t, err)
	j, err = m.JSON()
	assert.NoError(t, err)
	assert.Contains(t, j, "e_ZwRjAwRmZGtmZQR0ZQNjAwt2BGV5BN==")
	m, err = inbox.Fetch(13)
	assert.NoError(t, err)
	j, err = m.JSON()
	assert.NoError(t, err)
	assert.Contains(t, j, "e_ZwRjAwRmZGtmZwR0ZQNjAwt3AmxlZN==")
	m, err = inbox.Fetch(14)
	assert.NoError(t, err)
	j, err = m.JSON()
	assert.NoError(t, err)
	assert.Contains(t, j, "e_ZwRjAwRmZGtmZwN3ZQNjAwt3AmZlAD==")
}

func TestShrink(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_1.html",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_2.html",
		},
		{
			"GET",
			"https://yopmail.com/en/mail?b=test&id=me_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj%3D%3D",
			"features/mail.html",
		},
		{
			"GET",
			"https://yopmail.com/en/mail?b=test&id=me_ZwRjAwRmZGtmZGDlZQNjAwt3AGHkAt%3D%3D",
			"features/mail.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/consent?c=accept",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)

	err = inbox.ParseInboxPages(19)

	assert.NoError(t, err)
	assert.Equal(t, 19, inbox.Count())
	m, err := inbox.Fetch(0)
	assert.NoError(t, err)
	_, err = m.JSON()
	assert.NoError(t, err)
	_, err = inbox.Fetch(18)
	assert.NoError(t, err)
}

func TestShrinkEmptyInbox(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_empty.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/consent?c=accept",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)

	err = inbox.ParseInboxPages(1)

	assert.NoError(t, err)
	assert.Equal(t, 0, inbox.Count())
}

func TestShrinkWithLimitGreaterThanNumberOfMessagesAvailable(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_1.html",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_empty.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/consent?c=accept",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)

	err = inbox.ParseInboxPages(18)

	assert.NoError(t, err)
	assert.Equal(t, 15, inbox.Count())
}

func TestGetAll(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_1.html",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_2.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)

	err = inbox.ParseInboxPages(29)
	mails := inbox.InboxItems

	assert.NoError(t, err)
	assert.Len(t, mails, 29)
	assert.Equal(t, "e_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj==", mails[0].ID)
	assert.Equal(t, "e_ZwRjAwRmZGtmZQR0ZQNjAwt2BGV5BN==", mails[28].ID)
}

func TestFlush(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_1.html",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=2&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_2.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
		{
			"GET",
			"https://yopmail.com/consent?c=accept",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?login=test&p=1&d=all&ctrl=e_ZGtkZwVmZQNmBGV1ZQNjZQVjAwD1BD==&v=4.8&r_c=&id",
			"features/noop.html",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)

	err = inbox.ParseInboxPages(15)
	inbox.Flush()

	assert.NoError(t, err)
}

func TestFlushEmptyInbox(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	assert.NoError(t, registerResponders([]responder{
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_empty.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)

	err = inbox.ParseInboxPages(1)
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
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=&id=&login=test&p=1&r_c=&scrl=&spam=true&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/inbox_page_1.html",
		},
		{
			"GET",
			"https://yopmail.com",
			"features/main_page.html",
		},
		{
			"GET",
			"https://yopmail.com/ver/4.8/webmail.js",
			"features/webmail.js",
		},
		{
			"GET",
			"https://yopmail.com/en/inbox?ad=0&ctrl=&d=e_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj%3D%3D&id=&login=test&p=1&r_c=&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp",
			"features/noop.html",
		},
	}))

	inbox, err := NewInbox[client.MailHTMLDoc]("test", false)
	assert.NoError(t, err)

	err = inbox.ParseInboxPages(1)
	assert.NoError(t, inbox.Delete(0))

	assert.Equal(t, 1, httpmock.GetCallCountInfo()["GET https://yopmail.com/en/inbox?ad=0&ctrl=&d=e_ZwRjAwRmZGtmAwZ1ZQNjAwt5AQZmZj%3D%3D&id=&login=test&p=1&r_c=&v=4.8&yj=VZGV5AmpjZwp5ZGNmZwL0BQH&yp=UAQDkAGH2Amp2Zmt0ZmVmAGp"])
	assert.NoError(t, err)
}

func TestColoured(t *testing.T) {
	type scenario struct {
		name               string
		inbox              Inbox[client.MailHTMLDoc]
		outputExpected     string
		jsonOutputExpected string
		errorExpected      error
	}

	scenarios := []scenario{
		{
			name: "No mails in the inbox",
			inbox: Inbox[client.MailHTMLDoc]{
				Name:       "test",
				InboxItems: []InboxItem{},
			},
			errorExpected:      errors.New("inbox is empty"),
			jsonOutputExpected: `{"name":"test","mails":[]}`,
		},
		{
			name: "Display emails",
			inbox: Inbox[client.MailHTMLDoc]{
				Name: "test",
				InboxItems: []InboxItem{
					{
						ID:     "02d3583b-7b58-40cb-a2b7-c09d79673334",
						IsSPAM: true,
						Sender: &Sender{
							Mail: "test1@protonmail.com",
							Name: "test1",
						},
						Subject: "test1 subject",
					},
					{
						ID: "0343583b-7b58-40cb-a2b7-c09d79673334",
						Sender: &Sender{
							Mail: "test2@protonmail.com",
							Name: "test2",
						},
						Subject: "test2 subject",
					},
					{
						ID:     "0243583b-7b58-40cb-a2b7-c09d79673334",
						IsSPAM: true,
						Sender: &Sender{
							Mail: "test3@protonmail.com",
							Name: "test3",
						},
						Subject: "test3 subject",
					},
					{
						ID: "0783583b-7b58-40cb-a2b7-c09d79673334",
						Sender: &Sender{
							Name: "test4",
						},
						Subject: "test4 subject",
					},
					{
						ID: "0903583b-7b58-40cb-a2b7-c09d79673334",
						Sender: &Sender{
							Mail: "test5@protonmail.com",
						},
						Subject: "test5 subject",
					},
					{
						ID: "12d3583b-7b58-40cb-a2b7-c09d79673334",
						Sender: &Sender{
							Mail: "test6@protonmail.com",
							Name: "test6",
						},
					},
					{
						ID:      "67d3583b-7b58-40cb-a2b7-c09d79673334",
						Sender:  &Sender{},
						Subject: "test7 subject",
					},
					{
						ID:      "89d3583b-7b58-40cb-a2b7-c09d79673334",
						Subject: "test8 subject",
					},
					{
						ID:      "f44cf3b8-f6a4-4b75-b734-cb1553b23cf6",
						Subject: "test9 subject",
					},
					{
						ID:      "f207be30-fad5-4d73-aa30-f69cb2a5ebac",
						Subject: "test10 subject",
					},
					{
						ID:      "d64c2eeb-9ff6-4d33-b4dc-034557805308",
						Subject: "test11 subject",
					},
				},
			},
			outputExpected: ` 1 test1 <test1@protonmail.com> [SPAM]
   test1 subject

 2 test2 <test2@protonmail.com>
   test2 subject

 3 test3 <test3@protonmail.com> [SPAM]
   test3 subject

 4 test4
   test4 subject

 5 test5@protonmail.com
   test5 subject

 6 test6 <test6@protonmail.com>
   [no data to display]

 7 [no data to display]
   test7 subject

 8 [no data to display]
   test8 subject

 9 [no data to display]
   test9 subject

 10 [no data to display]
    test10 subject

 11 [no data to display]
    test11 subject`,
			jsonOutputExpected: `{"name":"test","mails":[{"id":"02d3583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test1@protonmail.com","name":"test1"},"subject":"test1 subject","isSPAM":true},{"id":"0343583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test2@protonmail.com","name":"test2"},"subject":"test2 subject","isSPAM":false},{"id":"0243583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test3@protonmail.com","name":"test3"},"subject":"test3 subject","isSPAM":true},{"id":"0783583b-7b58-40cb-a2b7-c09d79673334","sender":{"name":"test4"},"subject":"test4 subject","isSPAM":false},{"id":"0903583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test5@protonmail.com"},"subject":"test5 subject","isSPAM":false},{"id":"12d3583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test6@protonmail.com","name":"test6"},"subject":"","isSPAM":false},{"id":"67d3583b-7b58-40cb-a2b7-c09d79673334","sender":{},"subject":"test7 subject","isSPAM":false},{"id":"89d3583b-7b58-40cb-a2b7-c09d79673334","subject":"test8 subject","isSPAM":false},{"id":"f44cf3b8-f6a4-4b75-b734-cb1553b23cf6","subject":"test9 subject","isSPAM":false},{"id":"f207be30-fad5-4d73-aa30-f69cb2a5ebac","subject":"test10 subject","isSPAM":false},{"id":"d64c2eeb-9ff6-4d33-b4dc-034557805308","subject":"test11 subject","isSPAM":false}]}`,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			j, jerr := scenario.inbox.JSON()
			assert.NoError(t, jerr)
			assert.JSONEq(t, scenario.jsonOutputExpected, j)

			c, err := scenario.inbox.Coloured()
			if scenario.errorExpected != nil {
				assert.EqualError(t, err, scenario.errorExpected.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, scenario.outputExpected, c)
			}
		})
	}
}

type responder struct {
	method   string
	URL      string
	filename string
}

func registerResponders(responders []responder) error {
	for _, r := range responders {
		b, err := os.ReadFile(r.filename)
		if err != nil {
			return err
		}

		httpmock.RegisterResponder(r.method, r.URL,
			httpmock.NewStringResponder(200, string(b)))
	}
	return nil
}
