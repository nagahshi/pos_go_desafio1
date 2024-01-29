package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// NewConnDB - cria uma nova conexão com o banco de dados SQLite.
func NewConnDB(nameDB string) *sql.DB {
	db, err := sql.Open("sqlite3", nameDB)
	if err != nil {
		log.Fatal(err)
	}

	// Verifica a conexão com o banco de dados.
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Certifique-se de que a tabela não exista antes de criar
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bid TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
