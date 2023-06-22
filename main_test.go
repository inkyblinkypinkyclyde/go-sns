package main

import (
	"database/sql"
	"go-sns/email"
	"go-sns/models"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

var (
	calls     int
	events    []models.Event
	eventTime = sql.NullTime{Time: time.Now(), Valid: true}
	event     = models.Event{
		Inserted_at: eventTime,
		Ip_addr:     "1",
		Mac_addr:    "2",
		Subject:     "subject",
		Message:     "message",
	}
)

// for use in mocking only
type SpyMailSender struct {
}

func (s *SpyMailSender) SendMail(emailService *email.EmailService, emailConfig *email.EmailCreds, title, message string) {
	calls++
}

type SpyEventLogger struct {
}

func (s *SpyEventLogger) LogEvent(event models.Event) {
	events = append(events, event)
}

type SpyReceiver struct {
	Calls int
}

func (s *SpyReceiver) RecieveNewEventHttp() {
	s.Calls++
}

func teardown() {
	calls = 0
	events = nil
}

func TestRecieveEventAndLog(t *testing.T) {
	t.Run("should send mail", func(t *testing.T) {
		teardown()
		spyMailSender := &SpyMailSender{}
		spyMailSender.SendMail(emailService, emailConfig, "subject", "message")
		assert.Equal(t, calls, 1)

	})
	t.Run("should log an event", func(t *testing.T) {
		teardown()
		spyEventLogger := &SpyEventLogger{}

		spyEventLogger.LogEvent(event)
		assert.Equal(t, events[0], event)
	})
	t.Run("Should send mail and log event when request is recieved", func(y *testing.T) {
		teardown()
		newMailSender := &SpyMailSender{}
		newEventLogger := &SpyEventLogger{}

		receiver := &DefaultReceiver{
			EventLogger: newEventLogger,
			MailSender:  newMailSender,
			Event: models.EventRaw{
				Ip_addr:  "1",
				Mac_addr: "2",
				Subject:  "subject",
				Message:  "message",
			},
		}
		receiver.ProcessEvent(eventTime)
		assert.Equal(t, calls, 1)
		assert.Equal(t, events[0], event)
	})
}
