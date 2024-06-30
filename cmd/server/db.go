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

// initDB запускает скрипт создания БД.
func initDB(db *sql.DB) {
	list := strings.Split(initSQL, ";")
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			log.Printf("tx rollback error: %s", err)
		}
	}()

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
