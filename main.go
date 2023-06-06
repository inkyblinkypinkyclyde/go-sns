package main

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"time"

	events "go-sns/models"

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

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	logEvent()
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}
func main() {
	getEmailCreds()
	sendMail("This is just a test")
	connectDB()
	http.HandleFunc("/hello", getHello)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
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

func logEvent() {
	event := events.Event{
		Inserted_at: sql.NullTime{Time: time.Now(), Valid: true},
		Ip_addr:     "1.1.1.1",
		Mac_addr:    "00:00:00:00:00:00",
		Subject:     "Test",
		Message:     "This is a test",
	}
	_, err := db.NewInsert().Model(&event).Exec(context.Background())
	if err != nil {
		fmt.Printf("Error inserting event: %s\n", err)
	}
}
