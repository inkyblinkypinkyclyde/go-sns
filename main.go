package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	db "go-sns/db"
	"go-sns/email"
	"go-sns/ingestion"
	models "go-sns/models"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type EmailCreds struct {
	Address  string
	Password string
	SmtpHost string
	SmtpPort string
}

type IMailSender interface {
	SendMail(*email.EmailService, *email.EmailCreds, string, string)
}
type DefaultMailSender struct {
}

func (d *DefaultMailSender) SendMail(emailService *email.EmailService, emailConfig *email.EmailCreds, title, message string) {
	err := emailService.SendMail(email.ComposeEmail(emailConfig.Address, []string{emailConfig.Address}, title, message))
	if err != nil {
		fmt.Printf("Error sending email: %s\n", err)
	}
}

type IEventLogger interface {
	LogEvent(models.Event)
}
type DefaultEventLogger struct {
}

func (d *DefaultEventLogger) LogEvent(event models.Event) {
	_, err := myDb.NewInsert().Model(&event).Exec(context.Background())
	if err != nil {
		fmt.Printf("Error inserting event: %s\n", err)
	}
}

type DefaultReceiver struct {
	MailSender  IMailSender
	EventLogger IEventLogger
	Event       models.EventRaw
}

func (d *DefaultReceiver) GetMessageByCode(code string) { // TODO: test this function
	var messages []models.Message

	err := myDb.NewSelect().
		Model(&messages).
		Where("code = ?", code).
		Scan(context.Background())

	if err != nil {
		fmt.Printf("Error selecting code from db: %s", err)
	}
	d.Event.Subject = messages[0].Subject
	d.Event.Message = messages[0].Body
}

func (d *DefaultReceiver) ProcessEvent(now sql.NullTime) {
	if d.Event.Message == "bun" {
		d.GetMessageByCode(d.Event.Subject)
	}
	event := models.Event{
		Inserted_at: now,
		Ip_addr:     d.Event.Ip_addr,
		Mac_addr:    d.Event.Mac_addr,
		Subject:     d.Event.Subject,
		Message:     d.Event.Message,
	}
	messageBody := emailService.MessageBuilder(event)
	d.EventLogger.LogEvent(event)
	d.MailSender.SendMail(emailService, emailConfig, event.Subject, messageBody)
}

func (d *DefaultReceiver) RecieveNewEventHttp(c *gin.Context) { // TODO: test this function
	d.Event.Ip_addr = c.Param("ip_addr")
	d.Event.Mac_addr = c.Param("mac_addr")
	d.Event.Subject = c.Param("subject")
	d.Event.Message = c.Param("message")
	d.ProcessEvent(sql.NullTime{Time: time.Now(), Valid: true})
	d.ClearEvent()
}

func (d *DefaultReceiver) ClearEvent() {
	d.Event.Ip_addr = ""
	d.Event.Mac_addr = ""
	d.Event.Subject = ""
	d.Event.Message = ""
}

var (
	//go:embed emailcreds.json
	rawJson        string
	myDb           *bun.DB
	emailConfig, _ = email.GetConfig(rawJson)
	emailService   = email.NewEmailService(emailConfig)
)

func main() {
	myDb = db.ConnectDB()
	receiver := &DefaultReceiver{
		EventLogger: &DefaultEventLogger{},
		MailSender:  &DefaultMailSender{},
	}
	ingestion.IngestMessages(myDb)

	router := gin.Default()
	router.GET("http/:ip_addr/:mac_addr/:subject/:message", receiver.RecieveNewEventHttp)
	router.Run(":8080")
}
