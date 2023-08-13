package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSourceMail(t *testing.T) {
	type scenario struct {
		name               string
		mail               *SourceMail
		outputExpected     string
		jsonOutputExpected string
	}

	scenarios := []scenario{
		{
			name: "Display a regular email",
			mail: &SourceMail{
				ID: "e_ZwZjBQRmZwZkZwR4ZQNjAQZ4ZQRlZt==",
				Headers: map[string][]string{
					"Content-Transfer-Encoding": []string{"quoted-printable"},
					"Content-Type":              []string{"text/html; charset=utf-8"},
					"Date":                      []string{"Sun, 13 Aug 2023 23:12:14 +0000"},
					"Message-Id":                []string{"\u003c01000189f12beb37-090862b4-87cf-4a57-9071-9a39ade2308c-000000@email.amazonses.com\u003e"},
					"Mime-Version":              []string{"1.0"},
					"Sender":                    []string{"Ola no-reply \u003caplicativos@notificacionesatlas.com\u003e"},
					"Subject":                   []string{"=?utf-8?Q?Marcaci=C3=B3n?= de un punto de ronda fuera de la =?utf-8?Q?posici=C3=B3n?= georreferencia del cliente en PONTIFICIA UNIVERSIDAD JAVERIANA, zona:Casa Farallones."},
					"To":                        []string{"test@yopmail.com"},
					"X-Ses-Outgoing":            []string{"2023.08.13-54.240.48.112"},
				},
				Body: "TEST\nTEST\nTEST\n",
			},
			outputExpected: `---
Content-Transfer-Encoding : quoted-printable

Content-Type              : text/html; charset=utf-8

Date                      : Sun, 13 Aug 2023 23:12:14 +0000

Message-Id                : <01000189f12beb37-090862b4-87cf-4a57-9071-9a39ade230
                            8c-000000@email.amazonses.com>

Mime-Version              : 1.0

Sender                    : Ola no-reply <aplicativos@notificacionesatlas.com>

Subject                   : =?utf-8?Q?Marcaci=C3=B3n?= de un punto de ronda fuer
                            a de la =?utf-8?Q?posici=C3=B3n?= georreferencia del
                             cliente en PONTIFICIA UNIVERSIDAD JAVERIANA, zona:C
                            asa Farallones.

To                        : test@yopmail.com

X-Ses-Outgoing            : 2023.08.13-54.240.48.112

---
TEST
TEST
TEST

---
`,
			jsonOutputExpected: `{
"body":"TEST\nTEST\nTEST\n",
"headers":{
	"Content-Transfer-Encoding":["quoted-printable"],
	"Content-Type":["text/html; charset=utf-8"],
	"Date":["Sun, 13 Aug 2023 23:12:14 +0000"],
	"Message-Id":["<01000189f12beb37-090862b4-87cf-4a57-9071-9a39ade2308c-000000@email.amazonses.com>"],
	"Mime-Version":["1.0"],
	"Sender":["Ola no-reply <aplicativos@notificacionesatlas.com>"],
	"Subject":["=?utf-8?Q?Marcaci=C3=B3n?= de un punto de ronda fuera de la =?utf-8?Q?posici=C3=B3n?= georreferencia del cliente en PONTIFICIA UNIVERSIDAD JAVERIANA, zona:Casa Farallones."],
	"To":["test@yopmail.com"],
	"X-Ses-Outgoing":["2023.08.13-54.240.48.112"]
},
"id":"e_ZwZjBQRmZwZkZwR4ZQNjAQZ4ZQRlZt=="
}`,
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
