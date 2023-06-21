package main

import (
	"database/sql"
	"go-sns/models"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestRecieveEventAndLog(t *testing.T) {
	t.Run("should send mail", func(t *testing.T) {
		spyMailSender := &SpyMailSendOperations{}
		spyMailSender.SendMail()
		assert.Equal(t, spyMailSender.Calls, 1)
	})
	t.Run("should log an event", func(t *testing.T) {
		spyEventLogger := &SpyEventLoggingOperation{}
		eventTime := sql.NullTime{Time: time.Now(), Valid: true}
		event := models.Event{
			Inserted_at: eventTime,
			Ip_addr:     "1",
			Mac_addr:    "2",
			Subject:     "subject",
			Message:     "message",
		}
		spyEventLogger.LogEvent(event)
		assert.Equal(t, spyEventLogger.Events[0], event)
	})
}
