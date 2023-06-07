package main

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/smtp"
	"time"

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
	rawJson    string
	emailCreds EmailCreds
	db         *bun.DB
)

func main() {
	getEmailCreds()
	sendMail("go-sns is running")
	connectDB()

	router := gin.Default()
	router.GET(":ip_addr/:mac_addr/:subject/:message", recieveNewEvent)
	router.Run(":8080")
}

func getEmailCreds() {
	json.Unmarshal([]byte(rawJson), &emailCreds)
	fmt.Printf("Email address: %s\n", emailCreds.Address)
	fmt.Printf("Email password: %s\n", emailCreds.Password)
	fmt.Printf("SMTP Host is : %s\n", emailCreds.SmtpHost)
	fmt.Printf("SMTP Port is : %s\n", emailCreds.SmtpPort)
}

func sendMail(message string) {
	body := []byte("Subject: go-sns notification\r\n" + message)
	auth := smtp.PlainAuth("", emailCreds.Address, emailCreds.Password, emailCreds.SmtpHost)
	err := smtp.SendMail(emailCreds.SmtpHost+":"+emailCreds.SmtpPort, auth, emailCreds.Address, []string{emailCreds.Address}, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully")
}

func connectDB() {
	dsn := "postgresql://postgres:postgres@localhost:5434/sns-db?sslmode=disable"
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

func recieveNewEvent(c *gin.Context) {
	ip_addr := c.Param("ip_addr")
	mac_addr := c.Param("mac_addr")
	subject := c.Param("subject")
	message := c.Param("message")
	// fmt.Printf("column: %s, datum: %s\n", column, datum)
	event := events.Event{
		Inserted_at: sql.NullTime{Time: time.Now(), Valid: true},
		Ip_addr:     ip_addr,
		Mac_addr:    mac_addr,
		Subject:     subject,
		Message:     message,
	}
	logEvent(event)
	emailString := fmt.Sprintf("New event logged from\nIP: %s\nMAC: %s\nSubject: %s\nMessage: %s\n", ip_addr, mac_addr, subject, message)
	sendMail(emailString)
}
