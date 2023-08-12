package mail

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/jaytaylor/html2text"
)

const noDataToDisplayMsg = "[no data to display]"

// Sender defines a mail sender
type Sender struct {
	Mail string `json:"mail,omitempty"`
	Name string `json:"name,omitempty"`
}

// Mail is a mail message
type Mail struct {
	ID     string     `json:"id"`
	Sender *Sender    `json:"sender,omitempty"`
	Title  *string    `json:"title,omitempty"`
	Date   *time.Time `json:"date,omitempty"`
	Body   string     `json:"body,omitempty"`
	IsSPAM bool       `json:"isSPAM"`
}

func (m Mail) Coloured() (string, error) {
	info := struct {
		HasSenderName bool
		SenderName    string
		HasSenderMail bool
		SenderMail    string
		From          string
		Title         string
		Date          string
		Body          string
	}{}

	if m.Sender != nil {
		if m.Sender.Name != "" {
			info.SenderName = color.MagentaString(m.Sender.Name)
			info.HasSenderName = true
		} else {
			info.SenderName = color.MagentaString(noDataToDisplayMsg)
		}
		if m.Sender.Mail != "" {
			info.HasSenderMail = true
			info.SenderMail = color.MagentaString(m.Sender.Mail)
		} else {
			info.SenderMail = color.MagentaString(noDataToDisplayMsg)
		}
	} else {
		info.SenderName = color.MagentaString(noDataToDisplayMsg)
		info.SenderMail = color.MagentaString(noDataToDisplayMsg)
	}
	if m.Title != nil {
		info.Title = color.YellowString(*m.Title)
	} else {
		info.Title = color.YellowString(noDataToDisplayMsg)
	}
	if m.Date != nil {
		info.Date = color.GreenString(m.Date.Format("2006-01-02 15:04"))
	} else {
		info.Date = color.GreenString(noDataToDisplayMsg)
	}
	if m.Body != "" {
		info.Body = color.CyanString(m.Body)
	} else {
		info.Body = color.CyanString(noDataToDisplayMsg)
	}

	var buf bytes.Buffer
	tpl := template.Must(template.New("t").Parse(`---
From  : {{ if .HasSenderName -}}
{{- .SenderName -}}
{{- end -}}
{{- if (and .HasSenderMail .HasSenderName) }} {{ end -}}
{{- if (and (eq .HasSenderMail false) (eq .HasSenderName false)) }}{{ .SenderName }}{{- end -}}
{{- if .HasSenderMail -}}
	{{- if .HasSenderName -}}<{{- end -}}
	{{- .SenderMail -}}
	{{- if .HasSenderName -}}>{{- end -}}
{{ end }}
Title : {{.Title}}
Date  : {{.Date}}
---
{{.Body}}
---
`))
	err := tpl.Execute(&buf, info)
	return buf.String(), err
}

func (m Mail) JSON() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", errors.New("something wrong occurred")
	}
	s := string(data)
	return s, nil
}

func parseFrom(s string) (string, string) {
	re := regexp.MustCompile(`(?s)(.+?) <(.+?)>`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 3 {
		return strings.TrimSpace(matches[1]), matches[2]
	}

	re = regexp.MustCompile(`<(.+?)>`)
	matches = re.FindStringSubmatch(s)
	if len(matches) == 2 {
		return "", matches[1]
	}

	return "", ""
}

func parseDate(s string) *time.Time {
	date, err := time.Parse("Monday, January 02, 2006 3:04:05 PM", s)
	if err != nil {
		return nil
	}

	return &date
}

func parseHTML(content string, err error) string {
	if err != nil {
		return ""
	}
	text, err := html2text.FromString(content)
	if err != nil {
		return ""
	}

	return text
}

func Parse(doc *goquery.Document) Mail {
	mail := Mail{}
	doc.Find("body div.fl .ellipsis").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			title := strings.TrimSpace(s.Text())
			mail.Title = &title
		case 1:
			mail.Sender = &Sender{}
			mail.Sender.Name, mail.Sender.Mail = parseFrom(s.Text())
		case 2:
			mail.Date = parseDate(strings.Join(strings.Fields(s.Text()), " "))
		}
	})
	mail.Body = parseHTML(doc.Find("div#mail").Html())
	return mail
}
