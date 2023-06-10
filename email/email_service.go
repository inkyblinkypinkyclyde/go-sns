package email

import (
	_ "embed"
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
}

func NewEmailService(creds *EmailCreds) *EmailService {
	fmt.Println("NewEmailService")
	fmt.Println(creds)
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
