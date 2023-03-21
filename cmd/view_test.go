package cmd

import (
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
		outputExpected string
	}

	scenarios := []scenario{
		{
			name: "All infos defined",
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
		t.Run(scenario.name, func(t *testing.T) {
			current, err := computeMailOutput(scenario.mail)
			assert.NoError(t, err)
			assert.Equal(t, scenario.outputExpected, current)
		})
	}
}
