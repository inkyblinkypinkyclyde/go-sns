package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type EmailCreds struct {
	Address  string
	Password string
}

var (
	//go:embed emailcreds.json
	rawJson string
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
	var emailCreds EmailCreds
	json.Unmarshal([]byte(rawJson), &emailCreds)
	fmt.Printf("Email address: %s\n", emailCreds.Address)
	fmt.Printf("Email password: %s\n", emailCreds.Password)
	http.HandleFunc("/hello", getHello)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
