package mail

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/jaytaylor/html2text"
)

// Sender defines a mail sender
type Sender struct {
	Mail string `json:"mail,omitempty"`
	Name string `json:"name,omitempty"`
}

// HTMLMail is an HTML mail message
type HTMLMail struct {
	ID      string     `json:"id"`
	Sender  *Sender    `json:"sender,omitempty"`
	Subject string     `json:"subject,omitempty"`
	Date    *time.Time `json:"date,omitempty"`
	Body    string     `json:"body,omitempty"`
	IsSPAM  bool       `json:"isSPAM"`
}

func (m *HTMLMail) SetID(ID string) {
	m.ID = ID
}

func (m *HTMLMail) Coloured() (string, error) {
	info := struct {
		HasSenderName bool
		SenderName    string
		HasSenderMail bool
		SenderMail    string
		From          string
		Subject       string
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
	if m.Subject != "" {
		info.Subject = color.YellowString(m.Subject)
	} else {
		info.Subject = color.YellowString(noDataToDisplayMsg)
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
From    : {{ if .HasSenderName -}}
{{- .SenderName -}}
{{- end -}}
{{- if (and .HasSenderMail .HasSenderName) }} {{ end -}}
{{- if (and (eq .HasSenderMail false) (eq .HasSenderName false)) }}{{ .SenderName }}{{- end -}}
{{- if .HasSenderMail -}}
	{{- if .HasSenderName -}}<{{- end -}}
	{{- .SenderMail -}}
	{{- if .HasSenderName -}}>{{- end -}}
{{ end }}
Subject : {{.Subject}}
Date    : {{.Date}}
---
{{.Body}}
---
`))
	err := tpl.Execute(&buf, info)
	return buf.String(), err
}

func (m *HTMLMail) JSON() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", errors.New("something wrong occurred")
	}
	return string(data), nil
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
