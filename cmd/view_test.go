package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/antham/yogo/inbox"
	"github.com/stretchr/testify/assert"
)

func TestComputeMailOutput(t *testing.T) {
	date, err := time.Parse("2006-01-02 15:04", "2022-10-24 23:20")
	assert.NoError(t, err)

	type scenario struct {
		name           string
		mail           *inbox.Mail
		isJSONOutput   bool
		outputExpected string
	}

	scenarios := []scenario{
		{
			name: "Display a regular email",
			mail: &inbox.Mail{ID: "test", Sender: &inbox.Sender{Name: "test", Mail: "test@protonmail.com"}, Title: "A title", Date: &date, Body: "test"},
			outputExpected: `---
From  : test <test@protonmail.com>
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`,
		},
		{
			name:           "Display a regular email as JSON",
			isJSONOutput:   true,
			mail:           &inbox.Mail{ID: "test", Sender: &inbox.Sender{Name: "test", Mail: "test@protonmail.com"}, Title: "A title", Date: &date, Body: "test"},
			outputExpected: `{"id":"test","sender":{"mail":"test@protonmail.com","name":"test"},"title":"A title","date":"2022-10-24T23:20:00Z","body":"test","isSPAM":false}`,
		},
		{
			name: "No sender name defined",
			mail: &inbox.Mail{ID: "test", Sender: &inbox.Sender{Mail: "test@protonmail.com"}, Title: "A title", Date: &date, Body: "test"},
			outputExpected: `---
From  : test@protonmail.com
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`,
		},
		{
			name: "No sender email defined",
			mail: &inbox.Mail{ID: "test", Sender: &inbox.Sender{Name: "test"}, Title: "A title", Date: &date, Body: "test"},
			outputExpected: `---
From  : test
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`,
		},
		{
			name: "No sender informations defined",
			mail: &inbox.Mail{ID: "test", Sender: &inbox.Sender{}, Title: "A title", Date: &date, Body: "test"},
			outputExpected: `---
From  : [no data to display]
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`,
		},
		{
			name: "No sender object defined",
			mail: &inbox.Mail{ID: "test", Title: "A title", Date: &date, Body: "test"},
			outputExpected: `---
From  : [no data to display]
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`,
		},
		{
			name: "No title defined",
			mail: &inbox.Mail{ID: "test", Sender: &inbox.Sender{Name: "test", Mail: "test@protonmail.com"}, Date: &date, Body: "test"},
			outputExpected: `---
From  : test <test@protonmail.com>
Title : [no data to display]
Date  : 2022-10-24 23:20
---
test
---
`,
		},
		{
			name: "No date defined",
			mail: &inbox.Mail{ID: "test", Sender: &inbox.Sender{Name: "test", Mail: "test@protonmail.com"}, Title: "A title", Body: "test"},
			outputExpected: `---
From  : test <test@protonmail.com>
Title : A title
Date  : [no data to display]
---
test
---
`,
		},
		{
			name: "No body defined",
			mail: &inbox.Mail{ID: "test", Sender: &inbox.Sender{Name: "test", Mail: "test@protonmail.com"}, Title: "A title", Date: &date},
			outputExpected: `---
From  : test <test@protonmail.com>
Title : A title
Date  : 2022-10-24 23:20
---
[no data to display]
---
`,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			current, err := computeMailOutput(scenario.mail, scenario.isJSONOutput)
			assert.NoError(t, err)
			assert.Equal(t, scenario.outputExpected, current)
		})
	}
}

func TestComputeInboxMailOutput(t *testing.T) {
	type scenario struct {
		name           string
		inbox          inbox.Inbox
		isJSONOutput   bool
		outputExpected string
		errorExpected  error
	}

	scenarios := []scenario{
		{
			name:          "No mails in the inbox",
			errorExpected: errors.New("inbox is empty"),
		},
		{
			name: "Display emails",
			inbox: inbox.Inbox{
				Name: "test",
				Mails: []inbox.Mail{
					{
						IsSPAM: true,
						Sender: &inbox.Sender{
							Mail: "test1@protonmail.com",
							Name: "test1",
						},
						Title: "test1 title",
					},
					{
						Sender: &inbox.Sender{
							Mail: "test2@protonmail.com",
							Name: "test2",
						},
						Title: "test2 title",
					},
					{
						IsSPAM: true,
						Sender: &inbox.Sender{
							Mail: "test3@protonmail.com",
							Name: "test3",
						},
						Title: "test3 title",
					},
					{
						Sender: &inbox.Sender{
							Name: "test4",
						},
						Title: "test4 title",
					},
					{
						Sender: &inbox.Sender{
							Mail: "test5@protonmail.com",
						},
						Title: "test5 title",
					},
					{
						Sender: &inbox.Sender{
							Mail: "test6@protonmail.com",
							Name: "test6",
						},
					},
					{
						Sender: &inbox.Sender{},
						Title:  "test7 title",
					},
					{
						Title: "test8 title",
					},
				},
			},
			outputExpected: ` 1 test1 <test1@protonmail.com> [SPAM]
   test1 title

 2 test2 <test2@protonmail.com>
   test2 title

 3 test3 <test3@protonmail.com> [SPAM]
   test3 title

 4 test4
   test4 title

 5 test5@protonmail.com
   test5 title

 6 test6 <test6@protonmail.com>
   [no data to display]

 7 [no data to display]
   test7 title

 8 [no data to display]
   test8 title`,
		},
		{
			name: "Display emails as JSON",
			inbox: inbox.Inbox{
				Name: "test",
				Mails: []inbox.Mail{
					{
						ID:     "02d3583b-7b58-40cb-a2b7-c09d79673334",
						IsSPAM: true,
						Sender: &inbox.Sender{
							Mail: "test1@protonmail.com",
							Name: "test1",
						},
						Title: "test1 title",
					},
					{
						ID: "0343583b-7b58-40cb-a2b7-c09d79673334",
						Sender: &inbox.Sender{
							Mail: "test2@protonmail.com",
							Name: "test2",
						},
						Title: "test2 title",
					},
					{
						ID:     "0243583b-7b58-40cb-a2b7-c09d79673334",
						IsSPAM: true,
						Sender: &inbox.Sender{
							Mail: "test3@protonmail.com",
							Name: "test3",
						},
						Title: "test3 title",
					},
					{
						ID: "0783583b-7b58-40cb-a2b7-c09d79673334",
						Sender: &inbox.Sender{
							Name: "test4",
						},
						Title: "test4 title",
					},
					{
						ID: "0903583b-7b58-40cb-a2b7-c09d79673334",
						Sender: &inbox.Sender{
							Mail: "test5@protonmail.com",
						},
						Title: "test5 title",
					},
					{
						ID: "12d3583b-7b58-40cb-a2b7-c09d79673334",
						Sender: &inbox.Sender{
							Mail: "test6@protonmail.com",
							Name: "test6",
						},
					},
					{
						ID:     "67d3583b-7b58-40cb-a2b7-c09d79673334",
						Sender: &inbox.Sender{},
						Title:  "test7 title",
					},
					{
						ID:    "89d3583b-7b58-40cb-a2b7-c09d79673334",
						Title: "test8 title",
					},
				},
			},
			isJSONOutput:   true,
			outputExpected: `{"name":"test","mails":[{"id":"02d3583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test1@protonmail.com","name":"test1"},"title":"test1 title","isSPAM":true},{"id":"0343583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test2@protonmail.com","name":"test2"},"title":"test2 title","isSPAM":false},{"id":"0243583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test3@protonmail.com","name":"test3"},"title":"test3 title","isSPAM":true},{"id":"0783583b-7b58-40cb-a2b7-c09d79673334","sender":{"name":"test4"},"title":"test4 title","isSPAM":false},{"id":"0903583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test5@protonmail.com"},"title":"test5 title","isSPAM":false},{"id":"12d3583b-7b58-40cb-a2b7-c09d79673334","sender":{"mail":"test6@protonmail.com","name":"test6"},"title":"","isSPAM":false},{"id":"67d3583b-7b58-40cb-a2b7-c09d79673334","sender":{},"title":"test7 title","isSPAM":false},{"id":"89d3583b-7b58-40cb-a2b7-c09d79673334","title":"test8 title","isSPAM":false}]}`,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			current, err := computeInboxMailOutput(&scenario.inbox, scenario.isJSONOutput)
			if err != nil {
				assert.EqualError(t, err, scenario.errorExpected.Error())
			} else {
				assert.Equal(t, scenario.outputExpected, current)
			}
		})
	}
}
