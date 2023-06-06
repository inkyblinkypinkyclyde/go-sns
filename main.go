package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
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
	sendMail()
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

func sendMail() {
	message := []byte("This is just a test")
	auth := smtp.PlainAuth("", emailCreds.Address, emailCreds.Password, emailCreds.SmtpHost)
	err := smtp.SendMail(emailCreds.SmtpHost+":"+emailCreds.SmtpPort, auth, emailCreds.Address, []string{emailCreds.Address}, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully")
}
func logEvent() {}
