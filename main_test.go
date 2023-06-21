package main

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestSendMail(t *testing.T) {
	t.Run("should send mail", func(t *testing.T) {
		spyMailSender := &SpyMailSendOperations{}
		spyMailSender.SendMail()
		assert.Equal(t, spyMailSender.Calls, 1)

	})
}
