package models

type Message struct {
	Id      int    `bun:"id,pk,autoincrement"`
	Code    string `bun:"code"`
	Subject string `bun:"subject"`
	Body    string `bun:"body"`
}
