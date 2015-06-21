package mailbox

import "time"

type Mail struct {
	id         string
	FromString string
	FromMail   string
	SumUp      string
	Title      string
	at         time.Time
	headers    []string
	Body       string
}

func NewMail() {
}
