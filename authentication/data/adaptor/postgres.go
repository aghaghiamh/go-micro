package adaptor

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresDB struct {
	db *sql.DB
}

type Effector func() *PostgresDB

func OpenDB() *PostgresDB {
	// TODO: Add DSN (Data Source Name) to Config
	db, openErr := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if openErr != nil {
		log.Printf("couldn't connect to database: %s", openErr.Error())
		return nil
	}

	pErr := db.Ping()
	if pErr != nil {
		log.Printf("couldn't ping database: %s", pErr.Error())
		return nil
	}

	return &PostgresDB{
		db: db,
	}
}

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func() *PostgresDB {
		for r := 0; ; r++ {
			connectedDB := effector()

			if connectedDB != nil || r >= retries {
				return connectedDB
			}

			log.Printf("Attempt %d failed; Postgres is not yet ready; retrying in %v", r+1, delay)
			select {
			case <-time.After(delay):
				return nil
			}
		}
	}
}
