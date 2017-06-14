package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// POST /accounts
// GET /accounts/:id
// POST /accounts/:account_id/sites
// GET /accounts/:account_id/sites/:id

type H map[string]string

type Request struct {
	*http.Request

	body *[]byte
}

func NewRequest(r *http.Request) *Request {
	return &Request{r, nil}
}

func (r *Request) cachedBody() []byte {
	body := []byte{}

	if r.body == nil {
		tmp, err := ioutil.ReadAll(r.Body)

		if err == nil {
			body = tmp
		}

		r.body = &body
	}

	return *r.body
}

func (r *Request) UnmarshalTo(data interface{}) (err error) {
	err = json.Unmarshal(r.cachedBody(), &data)
	return
}

func (r *Request) MatchURL(pat string) (matches []string, ok bool) {
	re, _ := regexp.Compile(pat)
	m := re.FindStringSubmatch(r.URL.Path)

	if len(m) == 0 {
		ok = false
		matches = []string{}
	} else {
		ok = true
		matches = m[1:]
	}

	return
}

func (r *Request) MethodIn(verb string, verbs ...string) bool {
	all := append([]string{verb}, verbs...)

	for _, v := range all {
		if r.Method == v {
			return true
		}
	}

	return false
}

type Response struct {
	http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{w}
}

func (r *Response) JSON(code int, data interface{}) {
	b, _ := json.Marshal(data)

	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(code)
	r.Write(b)
	r.Write([]byte{'\n'})
}

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

		req := NewRequest(r)
		resp := NewResponse(w)

		if !req.MethodIn("POST") {
			resp.JSON(http.StatusBadRequest, H{"message": "Invalid request type"})
			return
		}

		err = req.UnmarshalTo(&account)

		if err != nil {
			resp.JSON(http.StatusBadRequest, H{"message": err.Error()})
			return
		}

		err = account.Create(db)

		if err != nil {
			resp.JSON(http.StatusBadRequest, H{"message": err.Error()})
			return
		}

		resp.JSON(http.StatusCreated, account)
	})

	mux.HandleFunc("/accounts/", func(w http.ResponseWriter, r *http.Request) {
		req := NewRequest(r)
		resp := NewResponse(w)

		if !req.MethodIn("GET") {
			resp.JSON(http.StatusBadRequest, H{"message": "Invalid request type"})
			return
		}

		matches, ok := req.MatchURL(`^/accounts/(.+)$`)

		if !ok {
			resp.JSON(http.StatusNotFound, H{"message": "Not found"})
			return
		}

		if account, _ := Find(db, matches[0]); account != nil {
			resp.JSON(http.StatusOK, account)
			return
		}

		resp.JSON(http.StatusNotFound, H{"message": "Not found"})
	})

	server := http.Server{Handler: mux, Addr: ":9000"}
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
