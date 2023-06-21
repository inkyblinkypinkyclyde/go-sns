package ingestion

import (
	"context"
	"fmt"
	"go-sns/models"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/uptrace/bun"
)

func IngestMessages(myDb *bun.DB) {
	TruncateTable(myDb, "messages")
	file := "ingestion/messages.xlsx"
	messages := IngestFile(file)
	for _, message := range messages {
		_, err := myDb.NewInsert().Model(&message).Exec(context.Background())
		if err != nil {
			fmt.Printf("Error inserting messages: %s\n", err)
		}
	}
}

func IngestFile(file string) []models.Message {
	f, err1 := excelize.OpenFile(file)
	if err1 != nil {
		fmt.Println(err1)
		fmt.Printf("Error opening file")
		return nil
	}

	firstSheet := f.WorkBook.Sheets.Sheet[0].Name

	var messages []models.Message

	rows := f.GetRows(firstSheet)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		messages = append(messages, models.Message{
			Code:    row[0],
			Subject: row[1],
			Body:    row[2],
		})

	}

	return messages
}

func TruncateTable(myDb *bun.DB, table string) {
	_, err := myDb.NewTruncateTable().Model(&models.Message{}).Exec(context.Background())
	if err != nil {
		fmt.Printf("Error truncating table: %s\n", err)
	}
}
