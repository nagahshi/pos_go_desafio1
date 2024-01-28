package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// NewConnDB - cria uma nova conexão com o banco de dados SQLite.
func NewConnDB(nameDB string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", nameDB)
	if err != nil {
		log.Fatal(err)
	}

	// Verifica a conexão com o banco de dados.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
