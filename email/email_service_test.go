package email

import (
	"testing"

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
