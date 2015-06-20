package mailbox

import "time"

type Mail struct {
	id      string
	from    string
	SumUp   string
	Title   string
	at      time.Time
	headers []string
	body    string
}

func NewMail() {

}
