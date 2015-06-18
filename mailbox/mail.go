package mailbox

import "time"

type Mail struct {
	title   string
	from    string
	at      time.Time
	headers []string
	body    string
}

func NewMail() {

}
