package postgresSql

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
)

type Postgres struct {
	db *sql.DB
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

type Data struct {
	IpAddress       string
	PingTime        string
	LastSuccessDate string
}

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

func (pg *Postgres) GetListIp() ([]Data, error) {
	const funcName = "GetListIp"
	str := "SELECT * FROM ip_addresses"
	rows, err := pg.db.Query(str)
	defer rows.Close()
	if err != nil {
		return []Data{}, fmt.Errorf("%s: %w", funcName, err)
	}

	var data []Data
	for rows.Next() {
		var row Data
		if err = rows.Scan(&row.IpAddress, &row.PingTime, &row.LastSuccessDate); err != nil {
			log.Println(err)
		}
		data = append(data, row)
	}
	return data, nil
}

func (pg *Postgres) NewEnty(data Data) error {
	const funcName = "newEnty"
	str := "INSERT INTO ip_addresses(ip_address, ping_time, last_success_date) VALUES ($1, $2, $3)"
	_, err := pg.db.Exec(str, data.IpAddress, data.PingTime, data.LastSuccessDate)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	return nil
}
