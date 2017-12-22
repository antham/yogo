package inbox

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestParseMailID(t *testing.T) {
	ID := parseMailID(`http://www.yopmail.com/en/mail.php?b=test&id=me_ZGLjAmN5ZGLmZwVlZQNjZmH3BQN5ZD==`)

	assert.Equal(t, ID, "me_ZGLjAmN5ZGLmZwVlZQNjZmH3BQN5ZD==", "Must extract mail ID")
}

func TestParseInboxPage(t *testing.T) {
	inbox := &Inbox{}

	parseInboxPage(getDoc("inbox_page"), inbox)

	assert.Equal(t, 15, inbox.Count(), "Must retrieve 15 mails")
	assert.Equal(t, "me_ZGLjAmN5ZGLmZmR5ZQNjZmH4ZQpmAt==", inbox.Get(0).ID, "Must retrieve mail ID")
	assert.Equal(t, "me_ZGLjAmN5ZGLlAmVkZQNjZmH2AQHlAD==", inbox.Get(14).ID, "Must retrieve mail ID")
}

func TestCount(t *testing.T) {
	inbox := &Inbox{}

	parseInboxPage(getDoc("inbox_page"), inbox)

	assert.Equal(t, inbox.Count(), 15, "Must retrieve 15 mails")
}

func TestParseInboxPages(t *testing.T) {
	fetchURL = func(URL string) (*goquery.Document, error) {
		URLS := map[string]string{
			"http://www.yopmail.com/en/inbox.php?login=test&p=1&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_page_1",
			"http://www.yopmail.com/en/inbox.php?login=test&p=2&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_page_2",
		}

		return getDoc(URLS[URL]), nil
	}

	inbox, err := ParseInboxPages("test", 29)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "test", inbox.GetIdentifier(), "Must return mailbox name")
	assert.Equal(t, 29, inbox.Count(), "Must return 30 elements")
	assert.Equal(t, "me_ZGLjAmN5ZGt1ZwVmZQNjAQNkZGZ3Aj==", inbox.Get(0).ID, "Must retrieve mail ID")
	assert.Equal(t, "me_ZGLjAmN5ZGt0ZmVlZQNjZmx3ZwHkZD==", inbox.Get(28).ID, "Must retrieve mail ID")
}

func TestShrink(t *testing.T) {
	fetchURL = func(URL string) (*goquery.Document, error) {
		URLS := map[string]string{
			"http://www.yopmail.com/en/inbox.php?login=test&p=1&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_page_1",
			"http://www.yopmail.com/en/inbox.php?login=test&p=2&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_page_2",
		}

		return getDoc(URLS[URL]), nil
	}

	inbox, err := ParseInboxPages("test", 19)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, 19, inbox.Count(), "Must return 19 elements")
	assert.Equal(t, "me_ZGLjAmN5ZGt1ZwVmZQNjAQNkZGZ3Aj==", inbox.Get(0).ID, "Must retrieve mail ID")
	assert.Equal(t, "me_ZGLjAmN5ZGt0AGH0ZQNjZmx4ZwHlZN==", inbox.Get(18).ID, "Must retrieve mail ID")
}

func TestShrinkEmptyInbox(t *testing.T) {
	fetchURL = func(URL string) (*goquery.Document, error) {
		URLS := map[string]string{
			"http://www.yopmail.com/en/inbox.php?login=test&p=1&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_empty",
		}

		return getDoc(URLS[URL]), nil
	}

	inbox, err := ParseInboxPages("test", 1)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, 0, inbox.Count(), "Must return 0 elements")
}

func TestShrinkWithLimitGreaterThanNumberOfMessagesAvailable(t *testing.T) {
	fetchURL = func(URL string) (*goquery.Document, error) {
		URLS := map[string]string{
			"http://www.yopmail.com/en/inbox.php?login=test&p=1&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_page_1",
			"http://www.yopmail.com/en/inbox.php?login=test&p=2&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_empty",
		}

		return getDoc(URLS[URL]), nil
	}

	inbox, err := ParseInboxPages("test", 18)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, 15, inbox.Count(), "Must return 15 elements")
}

func TestGetAll(t *testing.T) {
	fetchURL = func(URL string) (*goquery.Document, error) {
		URLS := map[string]string{
			"http://www.yopmail.com/en/inbox.php?login=test&p=1&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_page_1",
			"http://www.yopmail.com/en/inbox.php?login=test&p=2&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_page_2",
		}

		return getDoc(URLS[URL]), nil
	}

	inbox, err := ParseInboxPages("test", 29)
	mails := inbox.GetAll()

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, mails, 29, "Must return 29 elements")
	assert.Equal(t, "me_ZGLjAmN5ZGt1ZwVmZQNjAQNkZGZ3Aj==", mails[0].ID, "Must retrieve mail ID")
	assert.Equal(t, "me_ZGLjAmN5ZGt0ZmVlZQNjZmx3ZwHkZD==", mails[28].ID, "Must retrieve mail ID")
}

func TestFlush(t *testing.T) {
	fetchURL = func(URL string) (*goquery.Document, error) {
		URLS := map[string]string{
			"http://www.yopmail.com/en/inbox.php?login=test&p=1&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_page_1",
			"http://www.yopmail.com/en/inbox.php?login=test&p=2&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_page_2",
		}

		return getDoc(URLS[URL]), nil
	}

	send = func(URL string) {
		assert.Equal(t, "http://www.yopmail.com/en/inbox.php?login=test&p=1&d=all&ctrl=e_ZGLjAmN5ZGt1ZwVmZQNjAQNkZGZ3Aj==&v=2.7&r_c=&id=", URL, "Must build a correct deletion URL")
	}

	inbox, err := ParseInboxPages("test", 15)
	inbox.Flush()

	assert.NoError(t, err, "Must return no errors")
}

func TestFlushEmptyInbox(t *testing.T) {
	fetchURL = func(URL string) (*goquery.Document, error) {
		URLS := map[string]string{
			"http://www.yopmail.com/en/inbox.php?login=test&p=1&d=&ctrl=&scrl=&spam=true&v=2.7&r_c=&id=": "inbox_empty",
		}

		return getDoc(URLS[URL]), nil
	}

	inbox, err := ParseInboxPages("test", 1)
	inbox.Flush()

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, 0, inbox.Count(), "Must return 0 elements")
}
