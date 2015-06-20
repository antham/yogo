package mailbox

import "time"

type Mail struct {
	id      string
	sumUp   string
	title   string
	from    string
	at      time.Time
	headers []string
	body    string
}

func NewMail() {

}
