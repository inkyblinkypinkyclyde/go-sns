package email

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed emailcredstest.json
	jsonStringTest string
)

func TestGetConfig(t *testing.T) {
	emailCreds, err := GetConfig(jsonStringTest)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "username@domain.co.uk", emailCreds.Address)
	assert.Equal(t, "password", emailCreds.Password)
	assert.Equal(t, "smtp.domain.co.uk", emailCreds.SmtpHost)
	assert.Equal(t, 587, emailCreds.SmtpPort)
}
