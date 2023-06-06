package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"

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
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}
func main() {
	getEmailCreds()
	sendMail("This is just a test")
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
func logEvent(message string) {
	// ctx := context.Background()

	// Open a PostgreSQL database.
	dsn := "postgresql://postgres:postgres@localhost:5434/events"
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create a Bun db on top of it.
	db := bun.NewDB(pgdb, pgdialect.New())

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
}
