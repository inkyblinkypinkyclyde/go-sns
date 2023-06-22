package email

import (
	_ "embed"
	"fmt"
	events "go-sns/models"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
}

func NewEmailService(creds *EmailCreds) *EmailService {
	return &EmailService{
		dialer: gomail.NewDialer(creds.SmtpHost, creds.SmtpPort, creds.Address, creds.Password),
	}
}

func (e *EmailService) SendMail(email *gomail.Message) error {
	sender := email.GetHeader("From")[0]
	recipients := email.GetHeader("To")
	if err := e.dialer.DialAndSend(email); err != nil {
		return err
	}
	fmt.Printf("Email sent\n sender: %s\n recipients %s\n", sender, recipients[0])
	return nil
}

func (e *EmailService) MessageBuilder(event events.Event) string {
	eventTime := event.Inserted_at.Time.Format("2006-01-02 15:04:05")
	html := fmt.Sprintf(`<body>
    <p><b>A new event has been logged</b></p>
    <p>Time: %s</p>
    <p>From ip: %s</p>
    <p>From mac: %s</p>
    <p>Message: %s</p>
</body>`, eventTime, event.Ip_addr, event.Mac_addr, event.Message)
	return html
}
