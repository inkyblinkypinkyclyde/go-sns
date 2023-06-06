package events

import (
	"database/sql"
)

type Event struct {
	Id          int          `bun:"id,pk,autoincrement"`
	Inserted_at sql.NullTime `bun:"inserted_at"`
	Ip_addr     string       `bun:"ip_addr"`
	Mac_addr    string       `bun:"mac_addr"`
	Subject     string       `bun:"subject"`
	Message     string       `bun:"message"`
}
