package mailbox

import "time"

type Mail struct {
	Id         string
	FromString string
	FromMail   string
	SumUp      string
	Title      string
	Date       time.Time
	headers    []string
	Body       string
}
