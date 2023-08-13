package mail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"text/template"

	"github.com/fatih/color"
	"golang.org/x/term"
)

// SourceMail is an HTML  mail message
type SourceMail struct {
	ID      string              `json:"id"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func (m *SourceMail) SetID(ID string) {
	m.ID = ID
}

func (m *SourceMail) Coloured() (string, error) {
	var padding int
	info := struct {
		Headers map[string]string
		Body    string
	}{}
	width, _, err := term.GetSize(0)
	if err != nil {
		width = 80
	}
	info.Headers = map[string]string{}
	for k, vs := range m.Headers {
		for i, v := range vs {
			key := k
			if len(vs) > 1 {
				key = fmt.Sprintf("%s[%d]", k, i)
			}
			if len(key) > padding {
				padding = len(key)
			}
			info.Headers[key] = v
		}
	}
	info.Body = color.CyanString(m.Body)
	for k, v := range info.Headers {
		acc := []string{}
		for i, s := range splitString(v, width-3-padding) {
			if i == 0 {
				acc = append(acc, s)
				continue
			}
			acc = append(acc, fmt.Sprintf("%s%s", strings.Repeat(" ", 3+padding), s))
		}
		info.Headers[k] = color.MagentaString(strings.Join(acc, "\n"))
	}

	var buf bytes.Buffer
	tpl := template.Must(template.New("t").Parse(`---
{{ range $key, $value := .Headers -}}
{{ printf "%-` + strconv.Itoa(padding) + `s" $key }} : {{ $value }}

{{ end -}}
---
{{.Body}}
---
`))
	err = tpl.Execute(&buf, info)
	return buf.String(), err
}

func (m *SourceMail) JSON() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", errors.New("something wrong occurred")
	}
	return string(data), nil
}

func splitString(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	var b strings.Builder
	b.Grow(chunkSize)
	l := 0
	for _, r := range s {
		b.WriteRune(r)
		l++
		if l == chunkSize {
			chunks = append(chunks, b.String())
			l = 0
			b.Reset()
			b.Grow(chunkSize)
		}
	}
	if l > 0 {
		chunks = append(chunks, b.String())
	}
	return chunks
}
