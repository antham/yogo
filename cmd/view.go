package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"text/template"

	"github.com/fatih/color"

	"github.com/antham/yogo/inbox"
)

const noDataToDisplayMsg = "[no data to display]"

func computeInboxMailOutput(in Inbox, isJSONOutput bool) (string, error) {
	JSON, err := computeJSONOutput(in, isJSONOutput)
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
		info := struct {
			Index         string
			SenderName    string
			HasSenderName bool
			SenderMail    string
			HasSenderMail bool
			Title         string
			SPAM          string
		}{}

		if mail.Sender != nil {
			if mail.Sender.Name != "" {
				info.HasSenderName = true
				info.SenderName = color.YellowString(mail.Sender.Name)
			} else {
				info.SenderName = color.YellowString(noDataToDisplayMsg)
			}
			if mail.Sender.Mail != "" {
				info.HasSenderMail = true
				info.SenderMail = color.YellowString(mail.Sender.Mail)
			} else {
				info.SenderMail = color.YellowString(noDataToDisplayMsg)
			}
		} else {
			info.SenderName = color.YellowString(noDataToDisplayMsg)
			info.SenderMail = color.YellowString(noDataToDisplayMsg)
		}
		if mail.Title != "" {
			info.Title = color.CyanString(mail.Title)
		} else {
			info.Title = color.CyanString(noDataToDisplayMsg)
		}
		if mail.IsSPAM {
			info.SPAM = color.RedString("[SPAM]")
		}
		info.Index = strconv.Itoa(index + 1)

		var buf bytes.Buffer
		tpl := template.Must(template.New("t").Parse(` {{.Index}} {{ if .HasSenderName -}}
{{- .SenderName -}}
{{- end -}}
{{- if (and .HasSenderMail .HasSenderName) }} {{ end -}}
{{- if (and (eq .HasSenderMail false) (eq .HasSenderName false)) }}{{ .SenderName }}{{- end -}}
{{- if .HasSenderMail -}}
	{{- if .HasSenderName -}}<{{- end -}}
	{{- .SenderMail -}}
	{{- if .HasSenderName -}}>{{- end -}}
{{- end -}}
{{- if .SPAM }} {{ .SPAM -}}{{- end -}}
{{- if .Title }}
   {{ .Title }}
{{ end }}
`))
		if err := tpl.Execute(&buf, info); err != nil {
			return "", err
		}
		output = output + buf.String()
	}
	return strings.TrimRight(output, "\n"), nil
}

func computeMailOutput(mail *inbox.Mail, isJSONOutput bool) (string, error) {
	JSON, err := computeJSONOutput(*mail, isJSONOutput)
	if err != nil {
		return "", err
	}
	if JSON != nil {
		return *JSON, nil
	}

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

func computeJSONOutput(d interface{}, isJSONOutput bool) (*string, error) {
	if !isJSONOutput {
		return nil, nil
	}
	data, err := json.Marshal(d)
	if err != nil {
		return nil, errors.New("something wrong occurred")
	}
	s := string(data)
	return &s, nil
}

// info outputs a blue info message
func info(message string) string {
	return color.CyanString(message)
}

// success outputs a green successful message
func success(message string) string {
	return color.GreenString(message)
}
