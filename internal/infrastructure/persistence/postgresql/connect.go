package postgresql

import (
	"database/sql"
	"log"

	"github.com/IvanDrf/analyse-search-requests/internal/config"
	_ "github.com/lib/pq"
)

func Connect(config *config.PostgreSQLConfig) *sql.DB {
	db, err := sql.Open("postgres", config.DSN())
	if err != nil {
		log.Fatalf("can't connect to postgres database, error=%s", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("can't ping postgres database, error=%s", err)
	}

	return db
}
