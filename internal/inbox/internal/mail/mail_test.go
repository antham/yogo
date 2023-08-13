package mail

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antham/yogo/internal/client"
	"github.com/stretchr/testify/assert"
)

func getDoc[M client.MailDoc](t *testing.T, filename string) M {
	dir, err := os.Getwd()
	if err != nil {
		assert.NoError(t, err)
	}

	f, err := os.Open(dir + "/features/" + filename)
	if err != nil {
		assert.NoError(t, err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		assert.NoError(t, err)
	}

	return M(*doc)
}

func TestParseFrom(t *testing.T) {
	type scenario struct {
		name        string
		fromArg     string
		resultEmail string
		resultName  string
	}

	scenarios := []scenario{
		{
			name:        "parse name and email",
			fromArg:     "John Doe <john.doe@unknown.com>",
			resultName:  "John Doe",
			resultEmail: "john.doe@unknown.com",
		},
		{
			name: "parse name and email with spaces",
			fromArg: `Liana
                <AnnaMartinezpisea@lionspest.com.au>`,
			resultName:  "Liana",
			resultEmail: "AnnaMartinezpisea@lionspest.com.au",
		},
		{
			name:        "parse email only",
			fromArg:     "<john.doe@unknown.com>",
			resultName:  "",
			resultEmail: "john.doe@unknown.com",
		},
		{
			name:        "no email nor name to parse",
			fromArg:     "",
			resultName:  "",
			resultEmail: "",
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()
			name, mail := parseFrom(s.fromArg)
			assert.Equal(t, s.resultName, name)
			assert.Equal(t, s.resultEmail, mail)
		})
	}
}

func TestParseDate(t *testing.T) {
	date := parseDate("Sunday, June 13, 2021 8:57:08 PM")
	assert.Equal(t, "2021-06-13 20:57:08 +0000 UTC", date.String(), "Must parse date")

	date = parseDate("whatever")
	assert.Empty(t, date)
}

func TestParse(t *testing.T) {
	mail, err := Parse[client.MailHTMLDoc](getDoc[client.MailHTMLDoc](t, "mail.html"))
	assert.NoError(t, err)

	content, err := mail.Coloured()
	assert.NoError(t, err)

	assert.Equal(t, `---
From  : Liana <AnnaMartinezpisea@lionspest.com.au>
Title : In any case, I am happy that we met
Date  : 2021-06-13 20:57
---
( https://fectment.page.link/Ymry )

What such a gorgeous man is doing here?

*s ho Dent blink scorn league rose ivy superman atkins atkins mugsy freeze thorne katana bane jason edward batarang alfred rumor edward. w ph Maxie vale bartok selina hangman batman young hugo knight freeze batgirl ragman jason batmobile fairchild mister grayson ghul solomon the. ot Elongated czonk diamond bennett batmobile martha hatter snake bruce swamp strange blink creeper abattoir flash sinestro falcone harley bane ragdoll. o* ( https://fectment.page.link/CF1b )

( https://matering.page.link/bAmq )

Will you come to me on the weekend?

*s ho Todd aquaman bullock falcone jester chase croc doom swamp sinestro hangman fairchild nocturna hangman creeper hangman caird aquaman kane barrow. w p Clench chill green canary metallo face robin shrike hatter riddler gleeson justice rumor batarang kane lucius ragman fox grey batmobile. ho Night gleeson oswald cluemaster abattoir ragman gleeson oswald elongated batmobile face quinn abbott clayface moth knight prey knight atkins killer? to* ( https://exteleer.page.link/kjcS )
---
`, content)
}

func TestParseHTML(t *testing.T) {
	type scenario struct {
		name         string
		contentArg   string
		errorArg     error
		resultString string
	}

	scenarios := []scenario{
		{
			name:         "error provided is not nil",
			contentArg:   "",
			errorArg:     errors.New("an error occurred"),
			resultString: "",
		},
		{
			name:         "extract text from HTML",
			contentArg:   "<html>text</html>",
			errorArg:     nil,
			resultString: "text",
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()
			str := parseHTML(s.contentArg, s.errorArg)
			assert.Equal(t, s.resultString, str)
		})
	}
}

func TestMail(t *testing.T) {
	date, err := time.Parse("2006-01-02 15:04", "2022-10-24 23:20")
	assert.NoError(t, err)

	type scenario struct {
		name               string
		mail               *HTMLMail
		outputExpected     string
		jsonOutputExpected string
	}

	scenarios := []scenario{
		{
			name: "Display a regular email",
			mail: &HTMLMail{ID: "test", Sender: &Sender{Name: "test", Mail: "test@protonmail.com"}, Title: "A title", Date: &date, Body: "test"},
			outputExpected: `---
From  : test <test@protonmail.com>
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`,
			jsonOutputExpected: `{"id": "test", "isSPAM": false, "sender": {"name": "test", "mail": "test@protonmail.com"}, "title": "A title", "date": "2022-10-24T23:20:00Z", "body": "test"}`,
		},
		{
			name: "No sender name defined",
			mail: &HTMLMail{ID: "test", Sender: &Sender{Mail: "test@protonmail.com"}, Title: "A title", Date: &date, Body: "test"},
			outputExpected: `---
From  : test@protonmail.com
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`,
			jsonOutputExpected: `{"id":"test", "isSPAM": false, "sender": {"mail": "test@protonmail.com"}, "title": "A title", "date": "2022-10-24T23:20:00Z", "body": "test"}`,
		},
		{
			name: "No sender email defined",
			mail: &HTMLMail{ID: "test", Sender: &Sender{Name: "test"}, Title: "A title", Date: &date, Body: "test"},
			outputExpected: `---
From  : test
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`,
			jsonOutputExpected: `{"id":"test", "isSPAM": false, "sender": {"name": "test"}, "title": "A title", "date": "2022-10-24T23:20:00Z", "body": "test"}`,
		},
		{
			name: "No sender object defined",
			mail: &HTMLMail{ID: "test", Title: "A title", Date: &date, Body: "test"},
			outputExpected: `---
From  : [no data to display]
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`,
			jsonOutputExpected: `{"id":"test", "isSPAM": false, "title": "A title", "date": "2022-10-24T23:20:00Z", "body": "test"}`,
		},
		{
			name: "No title defined",
			mail: &HTMLMail{ID: "test", Sender: &Sender{Name: "test", Mail: "test@protonmail.com"}, Date: &date, Body: "test"},
			outputExpected: `---
From  : test <test@protonmail.com>
Title : [no data to display]
Date  : 2022-10-24 23:20
---
test
---
`,
			jsonOutputExpected: `{"id": "test", "isSPAM": false, "sender": {"name": "test", "mail": "test@protonmail.com"}, "date": "2022-10-24T23:20:00Z", "body": "test"}`,
		},
		{
			name: "No date defined",
			mail: &HTMLMail{ID: "test", Sender: &Sender{Name: "test", Mail: "test@protonmail.com"}, Title: "A title", Body: "test"},
			outputExpected: `---
From  : test <test@protonmail.com>
Title : A title
Date  : [no data to display]
---
test
---
`,
			jsonOutputExpected: `{"id":"test", "isSPAM": false, "sender": {"name": "test", "mail": "test@protonmail.com"}, "title": "A title", "body": "test"}`,
		},
		{
			name: "No body defined",
			mail: &HTMLMail{ID: "test", Sender: &Sender{Name: "test", Mail: "test@protonmail.com"}, Title: "A title", Date: &date},
			outputExpected: `---
From  : test <test@protonmail.com>
Title : A title
Date  : 2022-10-24 23:20
---
[no data to display]
---
`,
			jsonOutputExpected: `{"id": "test", "isSPAM": false, "sender": {"name": "test", "mail": "test@protonmail.com"}, "title": "A title", "date":"2022-10-24T23:20:00Z"}`,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			j, err := scenario.mail.JSON()
			assert.NoError(t, err)
			c, err := scenario.mail.Coloured()
			assert.NoError(t, err)

			assert.NoError(t, err)
			assert.Equal(t, scenario.outputExpected, c)
			assert.JSONEq(t, scenario.jsonOutputExpected, j)
		})
	}
}
