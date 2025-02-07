package postgresSql

import (
	"database/sql"
	"fmt"
	"sync"
)

type Postgres struct {
	db *sql.DB
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewPg(conn string) (*Postgres, error) {
	pgOnce.Do(func() {
		db, err := sql.Open("postgres", conn)
		if err != nil {
			fmt.Errorf("unable to create conection pool: %w", err)
		}
		pgInstance = &Postgres{db: db}
	})
	if pgInstance == nil {
		return nil, fmt.Errorf("failed to initialize Postgres instance")
	}
	return pgInstance, nil
}
