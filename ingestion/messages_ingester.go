package ingestion

import (
	"fmt"

	db "go-sns/db"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/uptrace/bun"
)

var (
	myDb *bun.DB
)

func IngestMessages() {
	f, err := excelize.OpenFile("messages.xlsx")
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error opening file")
		return
	}

	firstSheet := f.WorkBook.Sheets.Sheet[0].Name

	rows := f.GetRows(firstSheet)
	fmt.Println(rows)
	myDb := db.ConnectDB()
	myDb.NewInsert()
}
