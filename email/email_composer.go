package email

import "gopkg.in/gomail.v2"

func ComposeEmail(sender string, recievers []string, subject string, htmlBody string, attachments ...string) *gomail.Message {
	email := gomail.NewMessage()
	email.SetHeader("From", sender)
	email.SetHeader("To", recievers...)
	email.SetHeader("Subject", subject)
	email.SetBody("text/html", htmlBody)

	for _, attattachment := range attachments {
		email.Attach(attattachment)
	}
	return email
}
