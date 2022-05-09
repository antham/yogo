package inbox

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func getDoc(t *testing.T, filename string) *goquery.Document {
	dir, err := os.Getwd()
	if err != nil {
		assert.NoError(t, err)
	}

	f, err := os.Open(dir + "/" + filename)
	if err != nil {
		assert.NoError(t, err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		assert.NoError(t, err)
	}

	return doc
}

func TestParseFrom(t *testing.T) {
	name, mail := parseFrom(`John Doe <john.doe@unknown.com>`)
	assert.Equal(t, "John Doe", name)
	assert.Equal(t, "john.doe@unknown.com", mail)

	name, mail = parseFrom(`<john.doe@unknown.com>`)
	assert.Equal(t, "", name)
	assert.Equal(t, "john.doe@unknown.com", mail)

	name, mail = parseFrom(`Liana
                <AnnaMartinezpisea@lionspest.com.au>`)
	assert.Equal(t, "Liana", name)
	assert.Equal(t, "AnnaMartinezpisea@lionspest.com.au", mail)
}

func TestParseDate(t *testing.T) {
	date := parseDate("Sunday, June 13, 2021 8:57:08 PM")

	assert.Equal(t, "2021-06-13 20:57:08 +0000 UTC", date.String(), "Must parse date")
}

func TestParseMail(t *testing.T) {
	mail := &Mail{}
	parseMail(getDoc(t, "mail.html"), mail)

	assert.Equal(t, "Liana", mail.Sender.Name, "Must return sender name")
	assert.Equal(t, "AnnaMartinezpisea@lionspest.com.au", mail.Sender.Mail, "Must return sender email")
	assert.Equal(t, "In any case, I am happy that we met", mail.Title, "Must return mail title")
	assert.Equal(t, `( https://fectment.page.link/Ymry )

What such a gorgeous man is doing here?

*s ho Dent blink scorn league rose ivy superman atkins atkins mugsy freeze thorne katana bane jason edward batarang alfred rumor edward. w ph Maxie vale bartok selina hangman batman young hugo knight freeze batgirl ragman jason batmobile fairchild mister grayson ghul solomon the. ot Elongated czonk diamond bennett batmobile martha hatter snake bruce swamp strange blink creeper abattoir flash sinestro falcone harley bane ragdoll. o* ( https://fectment.page.link/CF1b )

( https://matering.page.link/bAmq )

Will you come to me on the weekend?

*s ho Todd aquaman bullock falcone jester chase croc doom swamp sinestro hangman fairchild nocturna hangman creeper hangman caird aquaman kane barrow. w p Clench chill green canary metallo face robin shrike hatter riddler gleeson justice rumor batarang kane lucius ragman fox grey batmobile. ho Night gleeson oswald cluemaster abattoir ragman gleeson oswald elongated batmobile face quinn abbott clayface moth knight prey knight atkins killer? to* ( https://exteleer.page.link/kjcS )`, mail.Body, "Must return mail body")
	assert.Equal(t, "2021-06-13 20:57:08 +0000 UTC", mail.Date.String(), "Must return mail date")
}
