package inbox

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func getDoc(filename string) *goquery.Document {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}

	f, err := os.Open(dir + "/" + filename)
	if err != nil {
		logrus.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		logrus.Fatal(err)
	}

	return doc
}

func TestParseFrom(t *testing.T) {
	name, mail := parseFrom(`De: "John Doe"   <john.doe@unknown.com>`)

	assert.Equal(t, "John Doe", name, "Must extract sender name")
	assert.Equal(t, "john.doe@unknown.com", mail, "Must extract sender email")
}

func TestParseDate(t *testing.T) {
	date := parseDate("Reçu le 07/07/2016 à 23:34")

	assert.Equal(t, "2016-07-07 23:34:00 +0000 UTC", date.String(), "Must parse date")
}

func TestParseMail(t *testing.T) {

	mail := &Mail{}
	parseMail(getDoc("mail.html"), mail)

	assert.Equal(t, "root", mail.Sender.Name, "Must return sender name")
	assert.Equal(t, "root@sd-50982.dedibox.fr", mail.Sender.Mail, "Must return sender email")
	assert.Equal(t, "Cronjob 123", mail.Title, "Must return mail title")
	assert.Equal(t, "Received: from root by sd-50982.dedibox.fr with local (Exim 4.87)\n(envelope-from )\nid 1bLnfc-0006p0-Cs\nfor test123@yopmail.com; Sat, 09 Jul 2016 10:31:28 +0200\nDate: Sat, 09 Jul 2016 10:31:28 +0200\nTo: test123@yopmail.com\nSubject: Cronjob 123\nUser-Agent: Heirloom mailx 12.5 7/5/10\nMIME-Version: 1.0\nContent-Type: text/plain; charset=us-ascii\nContent-Transfer-Encoding: quoted-printable\nMessage-Id:\nFrom: root\nX-AntiAbuse: This header was added to track abuse, please include it with any abuse report\nX-AntiAbuse: Primary Hostname - sd-50982.dedibox.fr\nX-AntiAbuse: Original Domain - yopmail.com\nX-AntiAbuse: Originator/Caller UID/GID - [0 0] / [47 12]\nX-AntiAbuse: Sender Address Domain - sd-50982.dedibox.fr\nX-Get-Message-Sender-Via: sd-50982.dedibox.fr: authenticated_id: root/primary_hostname/system user\nX-Authenticated-Sender: sd-50982.dedibox.fr: root\nX-Source: /usr/lib/systemd/systemd\nX-Source-Args: /usr/lib/systemd/systemd --system --deserialize 21\nX-Source-Dir: /root\n\n=", mail.Body, "Must return mail body")
	assert.Equal(t, "2016-07-09 10:31:00 +0000 UTC", mail.Date.String(), "Must return mail date")
}
