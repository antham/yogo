package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"text/template"

	"github.com/fatih/color"

	"github.com/antham/yogo/inbox"
)

var ErrSomethingWrongOccurred = errors.New("something wrong occurred")

func computeInboxMailOutput(in Inbox) (string, error) {
	JSON, err := computeJSONOutput(in)
	if err != nil {
		return "", err
	}
	if JSON != nil {
		return *JSON, nil
	}

	if in.Count() == 0 {
		return "", errors.New("inbox is empty")
	}

	output := ""
	for index, mail := range in.GetMails() {
		var spam string
		if mail.IsSPAM {
			spam = " [SPAM]"
		}
		output = output + fmt.Sprintf(" %s %s%s%s\n", color.GreenString(fmt.Sprintf("%d", index+1)), color.YellowString(mail.Sender.Mail), color.YellowString(mail.Sender.Name), color.RedString(spam))
		output = output + fmt.Sprintf(" %s\n\n", color.CyanString(mail.Title))
	}
	return output, nil
}

func computeMailOutput(mail *inbox.Mail) (string, error) {
	JSON, err := computeJSONOutput(*mail)
	if err != nil {
		return "", err
	}
	if JSON != nil {
		return *JSON, nil
	}

	const noDataToDisplayMsg = "[no data to display]"

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

	if mail.Sender != nil {
		if mail.Sender.Name != "" {
			info.SenderName = color.MagentaString(mail.Sender.Name)
			info.HasSenderName = true
		} else {
			info.SenderName = color.MagentaString(noDataToDisplayMsg)
		}
		if mail.Sender.Mail != "" {
			info.HasSenderMail = true
			info.SenderMail = color.MagentaString(mail.Sender.Mail)
		} else {
			info.SenderMail = color.MagentaString(noDataToDisplayMsg)
		}
	} else {
		info.SenderName = color.MagentaString(noDataToDisplayMsg)
		info.SenderMail = color.MagentaString(noDataToDisplayMsg)
	}
	if mail.Title != "" {
		info.Title = color.YellowString(mail.Title)
	} else {
		info.Title = color.YellowString(noDataToDisplayMsg)
	}
	if mail.Date != nil {
		info.Date = color.GreenString(mail.Date.Format("2006-01-02 15:04"))
	} else {
		info.Date = color.GreenString(noDataToDisplayMsg)
	}
	if mail.Body != "" {
		info.Body = color.CyanString(mail.Body)
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
	if err := tpl.Execute(&buf, info); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func computeJSONOutput(d interface{}) (*string, error) {
	if dumpJSON {
		data, err := json.Marshal(d)
		if err != nil {
			return nil, ErrSomethingWrongOccurred
		}
		s := string(data)
		return &s, nil
	}
	return nil, nil
}

// info outputs a blue info message
func info(message string) string {
	return color.CyanString(message)
}

// success outputs a green successful message
func success(message string) string {
	return color.GreenString(message)
}
