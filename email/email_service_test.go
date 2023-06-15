package email

import (
	"database/sql"
	events "go-sns/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEmailService(t *testing.T) {
	creds := &EmailCreds{
		Address:  "username@domain.co.uk",
		Password: "password",
		SmtpHost: "smtp.domain.co.uk",
		SmtpPort: 587,
	}
	emailService := NewEmailService(creds)
	assert.Equal(t, "smtp.domain.co.uk", emailService.dialer.Host)
	assert.Equal(t, 587, emailService.dialer.Port)
	assert.Equal(t, "username@domain.co.uk", emailService.dialer.Username)
	assert.Equal(t, "password", emailService.dialer.Password)
}

func TestMessageBuilder(t *testing.T) {
	event := &events.Event{
		Id:      1,
		Ip_addr: "192.168.1.1",
		Inserted_at: sql.NullTime{
			Time:  time.Date(2021, 01, 01, 00, 00, 00, 00, time.UTC),
			Valid: true,
		},
		Mac_addr: "00:00:00:00:00:00",
		Subject:  "Test Subject",
		Message:  "Test Message",
	}
	emailService := &EmailService{}
	message := emailService.MessageBuilder(*event)
	assert.Equal(t, "<body>\n    <p><b>A new event has been logged</b></p>\n    <p>Time: 2021-01-01 00:00:00</p>\n    <p>From ip: 192.168.1.1</p>\n    <p>From mac: 00:00:00:00:00:00</p>\n    <p>Message: Test Message</p>\n</body>", message)

}
