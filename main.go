package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// POST /accounts
// GET /accounts/:id
// POST /accounts/:account_id/sites
// GET /accounts/:account_id/sites/:id

type H map[string]string

type Account struct {
	ID   int64  `"json":"id"`
	Name string `"json":"name"`
	Slug string `"json":"slug"`
}

func (a *Account) Create(db *sql.DB) (err error) {
	stmt := `
		INSERT INTO accounts (name, slug)
		VALUES      (?, ?);
	`

	r, err := db.Exec(stmt, a.Name, a.Slug)

	if err != nil {
		return
	}

	a.ID, err = r.LastInsertId()

	if err != nil {
		return
	}

	return
}

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

	mux := http.NewServeMux()

	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		var (
			err     error
			account Account
		)

		w.Header().Add("Content-Type", "application/json")

		if r.Method != "POST" {
			b, _ := json.Marshal(H{"message": "Invalid request type"})

			w.WriteHeader(http.StatusBadRequest)
			w.Write(b)
			w.Write([]byte{'\n'})

			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &account)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = account.Create(db)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, _ := json.Marshal(account)

		w.Write(b)
		w.Write([]byte{'\n'})
	})

	server := http.Server{Handler: mux, Addr: ":9000"}
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
