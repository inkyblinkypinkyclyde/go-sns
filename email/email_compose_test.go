package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeEmail(t *testing.T) {
	sender := "username@domain.co.uk"
	recievers := []string{"alice@domain.co.uk", "bob@domain.co.uk"}
	subject := "Test Subject"
	htmlBody := "<h1>Test Body</h1>"
	attachments := []string{"test.txt", "test2.txt"}
	email := ComposeEmail(sender, recievers, subject, htmlBody, attachments...)
	assert.Equal(t, sender, email.GetHeader("From")[0])
	assert.Equal(t, recievers, email.GetHeader("To"))
	assert.Equal(t, subject, email.GetHeader("Subject")[0])
	// assert.Equal(t, htmlBody, email.getBody())}
}
