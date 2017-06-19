package main

import (
	"log"
	"net/http"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/reagent/copyright/app"
	"github.com/reagent/copyright/handlers"
)

// POST /accounts
// GET /accounts/:id
// POST /accounts/:account_id/sites
// GET /accounts/:account_id/sites/:id

func main() {
	var (
		db  *sql.DB
		err error
	)

	db, err = sql.Open("sqlite3", "./copyright.db")
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to database: %s\n", err.Error())
	}

	schema := `
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER NOT NULL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL
		);
	`

	if _, err = db.Exec(schema); err != nil {
		log.Fatalf("Schema creation failed: %s\n", err.Error())
	}

	a := app.New(db)

	a.POST(`^/accounts$`, handlers.AccountPost)
	a.GET(`^/accounts/(?P<id>\d+)$`, handlers.AccountGet)

	http.ListenAndServe(":9000", a)
}
