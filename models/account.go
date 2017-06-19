package models

import (
	"database/sql"
	"errors"
)

type Account struct {
	ID   int64  `"json":"id"`
	Name string `"json":"name"`
	Slug string `"json":"slug"`
}

func Find(db *sql.DB, id string) (*Account, error) {
	var count int

	err := db.QueryRow(
		"SELECT COUNT(*) FROM accounts WHERE id = ? OR slug = ?",
		id, id,
	).Scan(&count)

	if err != nil || count != 1 {
		return nil, errors.New("Record not found")
	}

	account := &Account{}

	stmt := `
		SELECT id, name, slug
		FROM   accounts
		WHERE  id = ? OR slug = ?;
	`

	err = db.QueryRow(stmt, id, id).Scan(&account.ID, &account.Name, &account.Slug)

	if err != nil {
		return nil, err
	}

	return account, nil
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
