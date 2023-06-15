package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"go-sns/email"
	events "go-sns/models"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type EmailCreds struct {
	Address  string
	Password string
	SmtpHost string
	SmtpPort string
}

var (
	//go:embed emailcreds.json
	rawJson        string
	db             *bun.DB
	emailConfig, _ = email.GetConfig(rawJson)
	emailService   = email.NewEmailService(emailConfig)
)

func main() {
	SendMail(emailService, emailConfig, "go-sns is running", "go-sns is running")

	connectDB()

	router := gin.Default()
	router.GET("http/:ip_addr/:mac_addr/:subject/:message", recieveNewEventHttp)
	router.Run(":8080")
}

func connectDB() {
	dsn := "postgresql://postgres:postgres@localhost:5435/sns-db?sslmode=disable"
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(pgdb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
}

func logEvent(event events.Event) {
	_, err := db.NewInsert().Model(&event).Exec(context.Background())
	if err != nil {
		fmt.Printf("Error inserting event: %s\n", err)
	}
}

func SendMail(emailService *email.EmailService, emailConfig *email.EmailCreds, title, message string) {
	err := emailService.SendMail(email.ComposeEmail(emailConfig.Address, []string{emailConfig.Address}, title, message))
	if err != nil {
		fmt.Printf("Error sending email: %s\n", err)
	}
}

func recieveNewEventHttp(c *gin.Context) {
	ip_addr := c.Param("ip_addr")
	mac_addr := c.Param("mac_addr")
	subject := c.Param("subject")
	message := c.Param("message")
	event := events.Event{
		Inserted_at: sql.NullTime{Time: time.Now(), Valid: true},
		Ip_addr:     ip_addr,
		Mac_addr:    mac_addr,
		Subject:     subject,
		Message:     message,
	}
	logEvent(event)
	messageBody := emailService.MessageBuilder(event)
	SendMail(emailService, emailConfig, subject, messageBody)

}
