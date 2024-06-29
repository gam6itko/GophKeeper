package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"strings"
)

//go:embed resources/init_db.sql
var initSQL string

// initDb запускает скрипт создания БД.
func initDb(db *sql.DB) {
	list := strings.Split(initSQL, ";")
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	for _, query := range list {
		query = strings.TrimSpace(query)
		if len(query) == 0 {
			continue
		}

		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
