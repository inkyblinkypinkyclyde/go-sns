package email

import (
	_ "embed"
	"encoding/json"
)

type EmailCreds struct {
	Address  string
	Password string
	SmtpHost string
	SmtpPort int
}

func GetConfig(jsonString string) (*EmailCreds, error) {
	var emailCreds EmailCreds
	if err := json.Unmarshal([]byte(jsonString), &emailCreds); err != nil {
		return nil, err
	}
	return &emailCreds, nil
}
