package mailbox

var baseUrl = "http://www.yopmail.com/en/inbox.php?login=%v&p=r&d=&ctrl=&scrl=&spam=true&v=2.6&r_c=&id="
type MailBox struct {
	mail  string
	mails []string
}

func NewMailBox() {

}

func (m *MailBox) Flush() bool {

}

func (m *MailBox) GetMails(emails string, limit int) bool {

}

func (m *MailBox) GetMail() bool {

}

func (m *Mail) Delete() bool {

}

func (m *Mail) GetAlias() string {

}
