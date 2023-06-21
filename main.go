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

// for use in mocking only
type SpyMailSendOperations struct {
	Calls int
}

func (s *SpyMailSendOperations) SendMail() {
	s.Calls++
}

type SpyEventLoggingOperation struct {
	Events []models.Event
}

func (s *SpyEventLoggingOperation) LogEvent(event models.Event) {
	s.Events = append(s.Events, event)
}

// for normal use
type DefaultMailSender struct {
}

func (d *DefaultMailSender) SendMail(emailService *email.EmailService, emailConfig *email.EmailCreds, title, message string) {
	err := emailService.SendMail(email.ComposeEmail(emailConfig.Address, []string{emailConfig.Address}, title, message))
	if err != nil {
		fmt.Printf("Error sending email: %s\n", err)
	}
}

type DefaultEventLogger struct {
}

func (d *DefaultEventLogger) LogEvent(event models.Event) {
	_, err := myDb.NewInsert().Model(&event).Exec(context.Background())
	if err != nil {
		fmt.Printf("Error inserting event: %s\n", err)
	}
}

var (
	//go:embed emailcreds.json
	rawJson        string
	myDb           *bun.DB
	emailConfig, _ = email.GetConfig(rawJson)
	emailService   = email.NewEmailService(emailConfig)
	mailSender     = &DefaultMailSender{}
	eventLogger    = &DefaultEventLogger{}
)

func main() {
	myDb = db.ConnectDB()
	ingestion.IngestMessages(myDb)

	router := gin.Default()
	router.GET("http/:ip_addr/:mac_addr/:subject/:message", RecieveNewEventHttp)
	router.Run(":8080")
}

func LogEvent(event models.Event) {

}

func RecieveNewEventHttp(c *gin.Context) {
	ip_addr := c.Param("ip_addr")
	mac_addr := c.Param("mac_addr")
	subject := c.Param("subject")
	message := c.Param("message")
	event := models.Event{
		Inserted_at: sql.NullTime{Time: time.Now(), Valid: true},
		Ip_addr:     ip_addr,
		Mac_addr:    mac_addr,
		Subject:     subject,
		Message:     message,
	}
	messageBody := emailService.MessageBuilder(event)
	eventLogger.LogEvent(event)
	mailSender.SendMail(emailService, emailConfig, subject, messageBody)
}
